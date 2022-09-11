package requests

type PasswordChangeRequest struct {
	BaseRequest `mapstructure:",squash"`
	Password    string `json:"password" mapstructure:"password"`
}
