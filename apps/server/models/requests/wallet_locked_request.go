package requests

type WalletLockedRequest struct {
	BaseRequest `mapstructure:",squash"`
	Wallet      string `json:"wallet" mapstructure:"wallet"`
}
