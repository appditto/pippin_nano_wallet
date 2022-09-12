package pow

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/appditto/pippin_nano_wallet/libs/pow/boompow"
	"github.com/appditto/pippin_nano_wallet/libs/pow/net"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/bbedward/nanopow"
	"k8s.io/klog/v2"
)

const BOOMPOW_URL = "https://boompow.banano.cc/graphql"

type PippinPow struct {
	WorkPeers        []string
	BpowClient       *boompow.BpowClient
	workPeersFailing bool
	mutex            sync.Mutex
}

func (p *PippinPow) WorkPeersFailing() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.workPeersFailing
}

func (p *PippinPow) SetWorkPeersFailing(failing bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.workPeersFailing = failing
}

// workPeers is an array of URLs to send work_generate requests to
func NewPippinPow(workPeers []string) *PippinPow {
	var bpowClient *boompow.BpowClient
	if utils.GetEnv("BPOW_KEY", "") != "" {
		bpowClient = boompow.NewBpowClient(utils.GetEnv("BPOW_URL", BOOMPOW_URL), utils.GetEnv("BPOW_KEY", ""))
	}
	return &PippinPow{
		BpowClient: bpowClient,
		WorkPeers:  workPeers,
		// If peers are failing we will generate local pow no matter what
		workPeersFailing: false,
	}
}

// Makes a request to configured array of work peers
func (p *PippinPow) workGenerateAPIRequest(ctx context.Context, url string, hash string, difficultyMultiplier int, difficulty string, validate bool, out chan *string) {
	fmt.Printf("Making work_generate request to %s\n", url)
	resp, err := net.MakeWorkGenerateRequest(ctx, url, hash, difficulty)
	if err == nil && resp.Work != "" {
		// Validate work
		if IsWorkValid(hash, difficultyMultiplier, resp.Work) || !validate {
			p.SetWorkPeersFailing(false)
			WriteChannelSafe(out, resp.Work)
		} else {
			klog.Errorf("Received invalid work %s for %s from %s", resp.Work, hash, url)
		}
	}
}

// Makes a request to BoomPoW
func (p *PippinPow) workGenerateBpowRequest(ctx context.Context, hash string, difficulty int, validate bool, out chan *string) {
	resp, err := p.BpowClient.WorkGenerate(ctx, hash, difficulty)
	if err == nil && resp != "" {
		if IsWorkValid(hash, difficulty, resp) || !validate {
			p.SetWorkPeersFailing(false)
			WriteChannelSafe(out, resp)
		} else {
			klog.Errorf("Received invalid work %s for %s from boompow", resp, hash)
		}
	}
}

// Use GPU or CPU to generate work
func (p *PippinPow) generateWorkLocally(hash string, difficultyMultiplier int) (string, error) {
	// Generate work locally
	if !utils.Validate64HexHash(hash) {
		return "", errors.New("invalid hash")
	}
	decoded, err := hex.DecodeString(hash)
	if err != nil {
		return "", err
	}
	res, err := nanopow.GenerateWork(decoded, DifficultyFromMultiplier(difficultyMultiplier))
	if err != nil {
		return "", err
	}
	return WorkToString(res), nil
}

// Will use OpenCL if compiled with -tags cl, otherwise pure golang implementation
func (p *PippinPow) workGenerateLocal(ctx context.Context, hash string, difficultyMultiplier int, validate bool, out chan *string) {
	// ! TODO - work out a way to cancel
	work, err := p.generateWorkLocally(hash, difficultyMultiplier)
	if err == nil {
		if IsWorkValid(hash, difficultyMultiplier, work) || !validate {
			WriteChannelSafe(out, work)
		} else {
			klog.Errorf("Received invalid work %s for %s from local", work, hash)
		}
	}
}

// Makes a work cancel request
func WorkCancelAPIRequest(url string, hash string) {
	net.MakeWorkCancelRequest(context.Background(), url, hash)
}

// Invokes work_generate requests to every peer simultaneously including BoomPoW, depending on configuration
// Returns the first valid work response, sends cancel to everybody else
// ! TODO - If no peers are configured, use a local work pool
// ! if peers are configured, use a local work peer only after failed requests
func (p *PippinPow) WorkGenerateMeta(hash string, difficultyMultiplier int, validate bool) (string, error) {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	resultChan := make(chan *string, len(p.WorkPeers))
	defer close(resultChan)

	difficultyUint := DifficultyFromMultiplier(difficultyMultiplier)
	difficultyStr := DifficultyToString(difficultyUint)
	runningLocally := false

	if (len(p.WorkPeers) < 1 && p.BpowClient == nil) || p.WorkPeersFailing() {
		// Local pow
		runningLocally = true
		go p.workGenerateLocal(ctx, hash, difficultyMultiplier, validate, resultChan)
	}
	for _, peer := range p.WorkPeers {
		go p.workGenerateAPIRequest(ctx, peer, hash, difficultyMultiplier, difficultyStr, validate, resultChan)
	}
	if p.BpowClient != nil {
		go p.workGenerateBpowRequest(ctx, hash, difficultyMultiplier, validate, resultChan)
	}

	select {
	case result := <-resultChan:
		// Send work cancel
		for _, peer := range p.WorkPeers {
			go WorkCancelAPIRequest(peer, hash)
		}
		return *result, nil
	// 30
	case <-time.After(10 * time.Second):
		// Send work cancel
		for _, peer := range p.WorkPeers {
			go WorkCancelAPIRequest(peer, hash)
		}
		// See if our peers are failing
		// Generate local pow if it didnt run locally
		if !runningLocally {
			p.SetWorkPeersFailing(true)
			work, err := p.generateWorkLocally(hash, difficultyMultiplier)
			if err == nil {
				return work, nil
			}
		}
		return "", errors.New("Unable to generate work")
	}
}

// Recovers from writing to close channel
func WriteChannelSafe(out chan *string, msg string) (err error) {
	defer func() {
		// recover from panic caused by writing to a closed channel
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			return
		}
	}()

	out <- &msg // write on possibly closed channel

	return err
}
