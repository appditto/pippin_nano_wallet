package requests

type AccountRequest struct {
	BaseRequest `mapstructure:",squash"`
	Account     string `json:"account" mapstructure:"account"`
}
