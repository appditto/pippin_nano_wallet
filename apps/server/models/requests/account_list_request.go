package requests

type AccountListRequest struct {
	BaseRequest `mapstructure:",squash"`
	Wallet      string       `json:"wallet" mapstructure:"wallet"`
	Count       *interface{} `json:"count,omitempty" mapstructure:"count,omitempty"`
}
