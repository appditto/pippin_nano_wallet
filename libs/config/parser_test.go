package config

import (
	"os"
	"path"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/config/models"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	os.Setenv("HOME", ".testdata")
	defer os.Unsetenv("HOME")
	defer os.RemoveAll(".testdata")
	configRoot, _ := utils.GetPippinConfigurationRoot()

	config, err := ParsePippinConfig()
	assert.Nil(t, err)

	// the config file should exist with defaults
	_, err = os.Stat(path.Join(configRoot, "config.yaml"))
	assert.Nil(t, err)

	// the config should have defaults
	assert.Equal(t, 11338, config.Server.Port)
	assert.Equal(t, "127.0.0.1", config.Server.Host)
	assert.Equal(t, "http://[::1]:7076", config.Server.NodeRpcUrl)
	assert.Equal(t, "", config.Server.NodeWsUrl)
	assert.Equal(t, false, config.Wallet.Banano)
	assert.Equal(t, true, *config.Wallet.AutoReceiveOnSend)
	assert.Equal(t, false, config.Wallet.NodeWorkGenerate)
	assert.Equal(t, []string{
		"ban_1ka1ium4pfue3uxtntqsrib8mumxgazsjf58gidh1xeo5te3whsq8z476goo",
		"ban_1cake36ua5aqcq1c5i3dg7k8xtosw7r9r7qbbf5j15sk75csp9okesz87nfn",
		"ban_1fomoz167m7o38gw4rzt7hz67oq6itejpt4yocrfywujbpatd711cjew8gjj",
	}, config.Wallet.PreconfiguredRepresentativesBanano)
	assert.Equal(t, []string{
		"nano_1x7biz69cem95oo7gxkrw6kzhfywq4x5dupw4z1bdzkb74dk9kpxwzjbdhhs",
		"nano_1thingspmippfngcrtk1ofd3uwftffnu4qu9xkauo9zkiuep6iknzci3jxa6",
		"nano_1natrium1o3z5519ifou7xii8crpxpk8y65qmkih8e8bpsjri651oza8imdd",
		"nano_3o7uzba8b9e1wqu5ziwpruteyrs3scyqr761x7ke6w1xctohxfh5du75qgaj",
	}, config.Wallet.PreconfiguredRepresentativesNano)
	assert.Equal(t, []string{}, config.Wallet.WorkPeers)
	assert.Equal(t, "1000000000000000000000000", config.Wallet.ReceiveMinimum)

	// Copy testdata config 1
	assert.Nil(t, os.Remove(path.Join(configRoot, "config.yaml")))
	// Read 1.yaml
	file, err := os.ReadFile(path.Join("testdata", "1.yaml"))
	assert.Nil(t, err)
	// Write 1.yaml to config.yaml
	assert.Nil(t, os.WriteFile(path.Join(configRoot, "config.yaml"), file, 0644))

	// ! Parse config - this one keeps all defaults except host
	config, err = ParsePippinConfig()
	assert.Nil(t, err)
	assert.Equal(t, 11338, config.Server.Port)
	assert.Equal(t, "0.0.0.0", config.Server.Host)
	assert.Equal(t, "http://[::1]:7076", config.Server.NodeRpcUrl)
	assert.Equal(t, "", config.Server.NodeWsUrl)
	assert.Equal(t, false, config.Wallet.Banano)
	assert.Equal(t, true, *config.Wallet.AutoReceiveOnSend)
	assert.Equal(t, false, config.Wallet.NodeWorkGenerate)
	assert.Equal(t, []string{
		"ban_1ka1ium4pfue3uxtntqsrib8mumxgazsjf58gidh1xeo5te3whsq8z476goo",
		"ban_1cake36ua5aqcq1c5i3dg7k8xtosw7r9r7qbbf5j15sk75csp9okesz87nfn",
		"ban_1fomoz167m7o38gw4rzt7hz67oq6itejpt4yocrfywujbpatd711cjew8gjj",
	}, config.Wallet.PreconfiguredRepresentativesBanano)
	assert.Equal(t, []string{
		"nano_1x7biz69cem95oo7gxkrw6kzhfywq4x5dupw4z1bdzkb74dk9kpxwzjbdhhs",
		"nano_1thingspmippfngcrtk1ofd3uwftffnu4qu9xkauo9zkiuep6iknzci3jxa6",
		"nano_1natrium1o3z5519ifou7xii8crpxpk8y65qmkih8e8bpsjri651oza8imdd",
		"nano_3o7uzba8b9e1wqu5ziwpruteyrs3scyqr761x7ke6w1xctohxfh5du75qgaj",
	}, config.Wallet.PreconfiguredRepresentativesNano)
	assert.Empty(t, config.Wallet.WorkPeers)
	assert.Equal(t, "1000000000000000000000000", config.Wallet.ReceiveMinimum)

	// Copy testdata config 2
	assert.Nil(t, os.Remove(path.Join(configRoot, "config.yaml")))
	// Read 1.yaml
	file, err = os.ReadFile(path.Join("testdata", "2.yaml"))
	assert.Nil(t, err)
	// Write 1.yaml to config.yaml
	assert.Nil(t, os.WriteFile(path.Join(configRoot, "config.yaml"), file, 0644))

	// ! Parse config - this one enabled banano which changes some defaults
	config, err = ParsePippinConfig()
	assert.Nil(t, err)
	assert.Equal(t, 11338, config.Server.Port)
	assert.Equal(t, "127.0.0.1", config.Server.Host)
	assert.Equal(t, "http://[::1]:7072", config.Server.NodeRpcUrl)
	assert.Equal(t, "", config.Server.NodeWsUrl)
	assert.Equal(t, true, config.Wallet.Banano)
	assert.Equal(t, true, *config.Wallet.AutoReceiveOnSend)
	assert.Equal(t, false, config.Wallet.NodeWorkGenerate)
	assert.Equal(t, []string{
		"ban_1ka1ium4pfue3uxtntqsrib8mumxgazsjf58gidh1xeo5te3whsq8z476goo",
		"ban_1cake36ua5aqcq1c5i3dg7k8xtosw7r9r7qbbf5j15sk75csp9okesz87nfn",
		"ban_1fomoz167m7o38gw4rzt7hz67oq6itejpt4yocrfywujbpatd711cjew8gjj",
	}, config.Wallet.PreconfiguredRepresentativesBanano)
	assert.Equal(t, []string{
		"nano_1x7biz69cem95oo7gxkrw6kzhfywq4x5dupw4z1bdzkb74dk9kpxwzjbdhhs",
		"nano_1thingspmippfngcrtk1ofd3uwftffnu4qu9xkauo9zkiuep6iknzci3jxa6",
		"nano_1natrium1o3z5519ifou7xii8crpxpk8y65qmkih8e8bpsjri651oza8imdd",
		"nano_3o7uzba8b9e1wqu5ziwpruteyrs3scyqr761x7ke6w1xctohxfh5du75qgaj",
	}, config.Wallet.PreconfiguredRepresentativesNano)
	assert.Empty(t, config.Wallet.WorkPeers)
	assert.Equal(t, "1000000000000000000000000000", config.Wallet.ReceiveMinimum)

	// Copy testdata config 3
	assert.Nil(t, os.Remove(path.Join(configRoot, "config.yaml")))
	// Read 3.yaml
	file, err = os.ReadFile(path.Join("testdata", "3.yaml"))
	assert.Nil(t, err)
	// Write 3.yaml to config.yaml
	assert.Nil(t, os.WriteFile(path.Join(configRoot, "config.yaml"), file, 0644))

	// ! Parse config - this one changes everything possible
	config, err = ParsePippinConfig()
	assert.Nil(t, err)
	assert.Equal(t, 500, config.Server.Port)
	assert.Equal(t, "1.2.3.4", config.Server.Host)
	assert.Equal(t, "https://coolnanonode.com/rpc", config.Server.NodeRpcUrl)
	assert.Equal(t, "ws://[::1]:7078", config.Server.NodeWsUrl)
	assert.Equal(t, true, config.Wallet.Banano)
	assert.Equal(t, false, *config.Wallet.AutoReceiveOnSend)
	assert.Equal(t, true, config.Wallet.NodeWorkGenerate)
	assert.Equal(t, []string{
		"ban_3tta9pdxr4djdcm6r3c7969syoirj3dunrtynmmi8n1qtxzk9iksoz1gxdrh",
	}, config.Wallet.PreconfiguredRepresentativesBanano)
	assert.Equal(t, []string{
		"nano_3tta9pdxr4djdcm6r3c7969syoirj3dunrtynmmi8n1qtxzk9iksoz1gxdrh",
	}, config.Wallet.PreconfiguredRepresentativesNano)
	assert.Equal(t, []string{
		"http://localhost:5555",
		"http://myotherworkpeer.com",
	}, config.Wallet.WorkPeers)
	assert.Equal(t, "1", config.Wallet.ReceiveMinimum)
}

func TestConfigValidation(t *testing.T) {
	os.Setenv("HOME", ".testdata")
	defer os.Unsetenv("HOME")
	defer os.RemoveAll(".testdata")
	// Setup
	configRoot := path.Join(".testdata", "config")
	os.RemoveAll(".testdata")
	os.MkdirAll(configRoot, 0755)

	// Copy testdata config 1
	// Read 1.yaml
	file, err := os.ReadFile(path.Join("testdata", "1.yaml"))
	assert.Nil(t, err)
	// Write 1.yaml to config.yaml
	assert.Nil(t, os.WriteFile(path.Join(configRoot, "config.yaml"), file, 0644))

	// ! Parse config - this one is valid
	config, err := ParsePippinConfig()
	assert.Nil(t, err)
	assert.Nil(t, config.Validate())

	// Set invalid rpc url
	config.Server.NodeRpcUrl = "httpz://[::1]:7072"
	assert.NotNil(t, config.Validate())
	assert.ErrorIs(t, config.Validate(), models.ErrInvalidRpcUrl)
	config.Server.NodeRpcUrl = "http://[::1]:7072"

	// Set invalid port
	config.Server.Port = 0
	assert.NotNil(t, config.Validate())
	assert.ErrorIs(t, config.Validate(), models.ErrInvalidPort)
	config.Server.Port = 11338

	// Check websocket port
	config.Server.NodeWsUrl = "ws://[::1]:7078"
	assert.Nil(t, config.Validate())
	config.Server.NodeWsUrl = "wsz://[::1]:7072"
	assert.NotNil(t, config.Validate())
	assert.ErrorIs(t, config.Validate(), models.ErrInvalidWSUrl)
	config.Server.NodeWsUrl = "ws://[::1]:7078"

	// Check receive minimum
	config.Wallet.ReceiveMinimum = "0"
	assert.NotNil(t, config.Validate())
	assert.ErrorIs(t, config.Validate(), models.ErrInvalidReceiveMinimum)
	config.Wallet.ReceiveMinimum = "133248290000000000000000000000000000001"
	assert.NotNil(t, config.Validate())
	assert.ErrorIs(t, config.Validate(), models.ErrInvalidReceiveMinimum)
	config.Wallet.ReceiveMinimum = "1"
	assert.Nil(t, config.Validate())

	// Check work peers
	config.Wallet.WorkPeers = []string{"http://localhost:5555", "http://myotherworkpeer.com"}
	assert.Nil(t, config.Validate())
	config.Wallet.WorkPeers = []string{"http://localhost:5555", "httpz://myotherworkpeer.com"}
	assert.NotNil(t, config.Validate())
	assert.ErrorContains(t, config.Validate(), "invalid work peer")
	config.Wallet.WorkPeers = []string{"http://localhost:5555", "http://myotherworkpeer.com"}

	// Check representatives
	config.Wallet.PreconfiguredRepresentativesBanano = []string{"ban_1fomoz167m7o38gw4rzt7hz67oq6itejpt4yocrfywujbpatd711cjew8gjj"}
	config.Wallet.PreconfiguredRepresentativesNano = []string{"nano_1fomoz167m7o38gw4rzt7hz67oq6itejpt4yocrfywujbpatd711cjew8gjj"}
	config.Wallet.Banano = false
	assert.Nil(t, config.Validate())
	config.Wallet.Banano = true
	assert.Nil(t, config.Validate())
	config.Wallet.PreconfiguredRepresentativesBanano = []string{"ban_1fomoz167m7o38gw4rzt7hz67oq6itejpt4yocrfywujbpatd711cjew8gjk"}
	assert.NotNil(t, config.Validate())
	assert.ErrorContains(t, config.Validate(), "invalid preconfigured representative")
	config.Wallet.Banano = false
	config.Wallet.PreconfiguredRepresentativesBanano = []string{}
	assert.Nil(t, config.Validate())
	config.Wallet.PreconfiguredRepresentativesNano = []string{"nano_1fomoz167m7o38gw4rzt7hz67oq6itejpt4yocrfywujbpatd711cjew8gjk"}
}
