package requests

type AccountsBalancesRequest struct {
	BaseRequest `mapstructure:",squash"`
	Accounts    []string `json:"accounts" mapstructure:"accounts"`
}
