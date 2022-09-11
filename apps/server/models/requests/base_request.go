package requests

// This is a common base for most requests
type BaseRequest struct {
	Action string `json:"action" mapstructure:"action"`
	Wallet string `json:"wallet" mapstructure:"wallet"`
}
