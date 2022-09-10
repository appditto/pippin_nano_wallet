package requests

type AccountListRequest struct {
	Action string       `json:"action" mapstructure:"action"`
	Wallet string       `json:"wallet" mapstructure:"wallet"`
	Count  *interface{} `json:"count,omitempty" mapstructure:"count,omitempty"`
}
