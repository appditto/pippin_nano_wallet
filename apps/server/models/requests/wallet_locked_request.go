package requests

type WalletLockedRequest struct {
	Action string `json:"action" mapstructure:"action"`
	Wallet string `json:"wallet" mapstructure:"wallet"`
}
