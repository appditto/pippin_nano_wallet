package requests

type AccountsRequest struct {
	BaseRequest `mapstructure:",squash"`
	Accounts    []string `json:"accounts" mapstructure:"accounts"`
}
