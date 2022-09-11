package requests

type PasswordEnterRequest struct {
	BaseRequest `mapstructure:",squash"`
	Password    string `json:"password" mapstructure:"password"`
}
