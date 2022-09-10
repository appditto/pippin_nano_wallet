package requests

type AccountCreateRequest struct {
	BaseRequest `mapstructure:",squash"`
	Wallet      string `json:"wallet" mapstructure:"wallet"`
}
