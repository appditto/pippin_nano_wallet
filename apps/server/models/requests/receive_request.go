package requests

type ReceiveRequest struct {
	BaseRequest `mapstructure:",squash"`
	Work        *string `json:"work,omitempty" mapstructure:"work,omitempty"`
	BpowKey     *string `json:"bpow_key,omitempty" mapstructure:"bpow_key,omitempty"`
	Account     string  `json:"account" mapstructure:"account"`
	Block       string  `json:"block" mapstructure:"block"`
}
