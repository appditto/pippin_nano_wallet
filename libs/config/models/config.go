package models

import (
	"errors"
	"math/big"
	"math/rand"
	"net/url"

	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"golang.org/x/exp/slices"
)

// ! The old server also had:
// log_file, log_to_stdout,
type ServerConfig struct {
	Host       string `yaml:"host" default:"127.0.0.1"`
	Port       int    `yaml:"port" default:"11338"`
	NodeRpcUrl string `yaml:"node_rpc_url"`
	NodeWsUrl  string `yaml:"node_ws_url"`
}

// ! The old server also had:
// max_work_processes, max_sign_threads
type WalletConfig struct {
	Banano                             bool     `yaml:"banano" default:"false"`
	PreconfiguredRepresentativesBanano []string `yaml:"preconfigured_representatives_banano" default:"[\"ban_1ka1ium4pfue3uxtntqsrib8mumxgazsjf58gidh1xeo5te3whsq8z476goo\",\"ban_1cake36ua5aqcq1c5i3dg7k8xtosw7r9r7qbbf5j15sk75csp9okesz87nfn\",\"ban_1fomoz167m7o38gw4rzt7hz67oq6itejpt4yocrfywujbpatd711cjew8gjj\"]"`
	PreconfiguredRepresentativesNano   []string `yaml:"preconfigured_representatives_nano" default:"[\"nano_1x7biz69cem95oo7gxkrw6kzhfywq4x5dupw4z1bdzkb74dk9kpxwzjbdhhs\",\"nano_1thingspmippfngcrtk1ofd3uwftffnu4qu9xkauo9zkiuep6iknzci3jxa6\",\"nano_1natrium1o3z5519ifou7xii8crpxpk8y65qmkih8e8bpsjri651oza8imdd\",\"nano_3o7uzba8b9e1wqu5ziwpruteyrs3scyqr761x7ke6w1xctohxfh5du75qgaj\"]"`
	WorkPeers                          []string `yaml:"work_peers"`
	NodeWorkGenerate                   bool     `yaml:"node_work_generate" default:"false"`
	ReceiveMinimum                     string   `yaml:"receive_minimum"`
	AutoReceiveOnSend                  *bool    `yaml:"auto_receive_on_send" default:"true"`
}

type PippinConfig struct {
	Server ServerConfig `yaml:"server"`
	Wallet WalletConfig `yaml:"wallet"`
}

// Implements the interface from creasty package
func (c *PippinConfig) SetDefaults() {
	if c.Wallet.Banano {
		if c.Wallet.ReceiveMinimum == "" {
			c.Wallet.ReceiveMinimum = "1000000000000000000000000000"
		}
		if c.Server.NodeRpcUrl == "" {
			c.Server.NodeRpcUrl = "http://[::1]:7072"
		}
	} else {
		if c.Wallet.ReceiveMinimum == "" {
			c.Wallet.ReceiveMinimum = "1000000000000000000000000"
		}
		if c.Server.NodeRpcUrl == "" {
			c.Server.NodeRpcUrl = "http://[::1]:7076"
		}
	}
}

var ErrInvalidRpcUrl = errors.New("invalid node_rpc_url")
var ErrInvalidWSUrl = errors.New("invalid node_ws_url")
var ErrInvalidPort = errors.New("invalid server port, out of range")
var ErrInvalidReceiveMinimum = errors.New("invalid receive_minimum, must be between 1 and 133248290000000000000000000000000000000 (max supply)")

func (c *PippinConfig) Validate() error {
	u, err := url.Parse(c.Server.NodeRpcUrl)
	if err != nil || !slices.Contains([]string{"http", "https"}, u.Scheme) || u.Host == "" {
		return ErrInvalidRpcUrl
	}

	// Parse server port as int
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return ErrInvalidPort
	}

	// Validate websocket URL if set
	if c.Server.NodeWsUrl != "" {
		u, err := url.Parse(c.Server.NodeWsUrl)
		if err != nil || !slices.Contains([]string{"ws", "wss"}, u.Scheme) || u.Host == "" {
			return ErrInvalidWSUrl
		}
	}

	// Parse receive minimum as big int
	minimum, ok := big.NewInt(0).SetString(c.Wallet.ReceiveMinimum, 10)
	if !ok {
		return ErrInvalidReceiveMinimum
	} else {
		maxSupply, _ := big.NewInt(0).SetString("133248290000000000000000000000000000000", 10)
		if minimum.Cmp(big.NewInt(1)) < 0 || minimum.Cmp(maxSupply) > 0 {
			return ErrInvalidReceiveMinimum
		}
	}

	// Validate all work peers
	for _, peer := range c.Wallet.WorkPeers {
		u, err := url.Parse(peer)
		if err != nil || !slices.Contains([]string{"http", "https"}, u.Scheme) || u.Host == "" {
			return errors.New("invalid work peer: " + peer)
		}
	}

	// Validate representatives
	if c.Wallet.Banano {
		for _, rep := range c.Wallet.PreconfiguredRepresentativesBanano {
			if _, err := utils.AddressToPub(rep); err != nil {
				return errors.New("invalid preconfigured representative: " + rep)
			}
		}
	} else {
		for _, rep := range c.Wallet.PreconfiguredRepresentativesBanano {
			if _, err := utils.AddressToPub(rep); err != nil {
				return errors.New("invalid preconfigured representative: " + rep)
			}
		}
	}

	return err
}

var ErrNoRepsConfigured = errors.New("no representatives configured")

func (c *PippinConfig) GetRandomRep() (string, error) {
	// Retrieve a random banano or nano representative from the arrays
	if c.Wallet.Banano {
		if len(c.Wallet.PreconfiguredRepresentativesBanano) == 0 {
			return "", ErrNoRepsConfigured
		}
		return c.Wallet.PreconfiguredRepresentativesBanano[rand.Intn(len(c.Wallet.PreconfiguredRepresentativesBanano))], nil
	}
	if len(c.Wallet.PreconfiguredRepresentativesNano) == 0 {
		return "", ErrNoRepsConfigured
	}
	return c.Wallet.PreconfiguredRepresentativesNano[rand.Intn(len(c.Wallet.PreconfiguredRepresentativesNano))], nil

}
