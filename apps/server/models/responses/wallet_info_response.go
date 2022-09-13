package responses

type WalletInfoResponse struct {
	Balance            string `json:"balance" mapstructure:"balance"`
	Pending            string `json:"pending" mapstructure:"pending"`
	Receivable         string `json:"receivable" mapstructure:"receivable"`
	AccountsCount      int    `json:"accounts_count" mapstructure:"accounts_count"`
	AdhocCount         int    `json:"adhoc_count" mapstructure:"adhoc_count"`
	DeterministicCount int    `json:"deterministic_count" mapstructure:"deterministic_count"`
	DeterministicIndex int    `json:"deterministic_index" mapstructure:"deterministic_index"`
}
