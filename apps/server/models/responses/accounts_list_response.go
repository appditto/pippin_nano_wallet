package responses

type AccountsListResponse struct {
	Accounts []string `json:"accounts" mapstructure:"accounts"`
}
