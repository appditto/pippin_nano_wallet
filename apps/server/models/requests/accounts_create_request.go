package requests

type AccountsCreateRequest struct {
	AccountCreateRequest `mapstructure:",squash"`
	Count                int `json:"count" mapstructure:"count"`
}
