package requests

// This is a common base for most requests
type BaseRequest struct {
	Action  string  `json:"action" mapstructure:"action"`
	Wallet  string  `json:"wallet" mapstructure:"wallet"`
	BpowKey *string `json:"bpow_key,omitempty" mapstructure:"bpow_key,omitempty"`
}
