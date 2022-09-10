package requests

type PasswordEnterRequest struct {
	Action   string `json:"action" mapstructure:"action"`
	Wallet   string `json:"wallet" mapstructure:"wallet"`
	Password string `json:"password" mapstructure:"password"`
}
