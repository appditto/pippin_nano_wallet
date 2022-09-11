package responses

type AccountBalanceItem struct {
	Balance    string `json:"balance" mapstructure:"balance"`
	Pending    string `json:"pending" mapstructure:"pending"`
	Receivable string `json:"receivable" mapstructure:"receivable"`
}

type AccountsBalancesResponse struct {
	Balances *map[string]AccountBalanceItem `json:"balances,omitempty" mapstructure:"balances,omitempty"`
}
