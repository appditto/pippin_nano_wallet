package models

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
