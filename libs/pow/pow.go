package pow

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/appditto/pippin_nano_wallet/libs/pow/net"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/bbedward/nanopow"
	"k8s.io/klog/v2"
)

type PippinPow struct {
	WorkPeers        []string
	workPeersFailing bool
	bpowKey          string
	bpowUrl          string
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
// bpowKey and bpowUrl are optional, bpowUrl will default to boompow.banano.cc/graphql
func NewPippinPow(workPeers []string, bpowKey string, bpowUrl string) *PippinPow {
	if bpowUrl == "" {
		bpowUrl = "https://boompow.banano.cc/graphql"
	}
	return &PippinPow{
		WorkPeers: workPeers,
		// If peers are failing we will generate local pow no matter what
		workPeersFailing: false,
		bpowUrl:          bpowUrl,
		bpowKey:          bpowKey,
	}
}

// Makes a request to configured array of work peers
func (p *PippinPow) workGenerateAPIRequest(ctx context.Context, url string, hash string, difficultyMultiplier int, difficulty string, validate bool, out chan *string) {
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
func (p *PippinPow) workGenerateBpowRequest(ctx context.Context, hash string, difficulty int, validate bool, blockAward bool, bpowKey string, out chan *string) {
	resp, err := net.MakeBoompowWorkGenerateRequest(ctx, p.bpowUrl, bpowKey, hash, difficulty, blockAward)
	if err == nil && resp != "" {
		// Validate work
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

// The main entry point for Pippin WorkGenerate
// Invokes work_generate requests to every peer simultaneously including BoomPoW, depending on configuration
// Returns the first valid work response, sends cancel to everybody else
// If no peers or boompow configured, uses local PoW
// If all peers fail, will use local PoW until peers are responsive again
func (p *PippinPow) WorkGenerateMeta(hash string, difficultyMultiplier int, validate bool, blockAward bool, bpowKey string) (string, error) {

	// 1 hard coded valid work is just for higher level integration tests so we don't need to calculate real work
	if hash == "3F93C5CD2E314FA16702189041E68E68C07B27961BF37F0B7705145BEFBA3AA3" {
		return "205452237a9b01f4", nil
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chanSize := len(p.WorkPeers)
	if p.bpowUrl != "" && (bpowKey != "" || p.bpowKey != "") {
		chanSize++
	}
	resultChan := make(chan *string, chanSize)
	defer close(resultChan)

	difficultyUint := DifficultyFromMultiplier(difficultyMultiplier)
	difficultyStr := DifficultyToString(difficultyUint)
	runningLocally := false

	if (len(p.WorkPeers) < 1 && p.bpowKey == "" && bpowKey == "") || p.WorkPeersFailing() {
		// Local pow
		runningLocally = true
		go p.workGenerateLocal(ctx, hash, difficultyMultiplier, validate, resultChan)
	}
	for _, peer := range p.WorkPeers {
		go p.workGenerateAPIRequest(ctx, peer, hash, difficultyMultiplier, difficultyStr, validate, resultChan)
	}
	if p.bpowUrl != "" {
		key := bpowKey
		if key == "" && p.bpowKey != "" {
			key = p.bpowKey
		}
		if key != "" {
			go p.workGenerateBpowRequest(ctx, hash, difficultyMultiplier, validate, blockAward, key, resultChan)
		}
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
