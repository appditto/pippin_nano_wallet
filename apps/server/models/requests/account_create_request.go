package requests

type AccountCreateRequest struct {
	Action string `json:"action" mapstructure:"action"`
	Wallet string `json:"wallet" mapstructure:"wallet"`
}
