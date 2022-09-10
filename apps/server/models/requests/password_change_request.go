package requests

type PasswordChangeRequest struct {
	Action   string `json:"action" mapstructure:"action"`
	Wallet   string `json:"wallet" mapstructure:"wallet"`
	Password string `json:"password" mapstructure:"password"`
}
