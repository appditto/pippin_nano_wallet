package requests

type WalletLockRequest struct {
	BaseRequest `mapstructure:",squash"`
	Wallet      string `json:"wallet" mapstructure:"wallet"`
}
