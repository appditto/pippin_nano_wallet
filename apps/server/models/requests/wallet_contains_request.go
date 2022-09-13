package requests

type WalletContainsRequest struct {
	BaseRequest `mapstructure:",squash"`
	Account     string `json:"account" mapstructure:"account"`
}
