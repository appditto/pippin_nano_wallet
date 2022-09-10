package responses

type AccountsCreateResponse struct {
	Accounts []string `json:"accounts" mapstructure:"accounts"`
}
