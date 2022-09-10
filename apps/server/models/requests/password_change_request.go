package requests

type PasswordChangeRequest struct {
	BaseRequest `mapstructure:",squash"`
	Wallet      string `json:"wallet" mapstructure:"wallet"`
	Password    string `json:"password" mapstructure:"password"`
}
