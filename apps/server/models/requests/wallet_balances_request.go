package requests

type WalletBalancesRequest struct {
	BaseRequest `mapstructure:",squash"`
	Wallet      string `json:"wallet" mapstructure:"wallet"`
}
