package requests

type WalletAddRequest struct {
	Action string `json:"action" mapstructure:"action"`
	Key    string `json:"key" mapstructure:"key"`
	Wallet string `json:"wallet" mapstructure:"wallet"`
}
