package responses

type AccountsResponse struct {
	Accounts []string `json:"accounts" mapstructure:"accounts"`
}
