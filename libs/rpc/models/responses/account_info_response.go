package responses

type AccountInfoResponse struct {
	Frontier                   string `json:"frontier" mapstructure:"frontier"`
	OpenBlock                  string `json:"open_block" mapstructure:"open_block"`
	RepresentativeBlock        string `json:"representative_block" mapstructure:"representative_block"`
	Balance                    string `json:"balance" mapstructure:"balance"`
	ModifiedTimestamp          string `json:"modified_timestamp" mapstructure:"modified_timestamp"`
	BlockCount                 string `json:"block_count" mapstructure:"block_count"`
	AccountVersion             string `json:"account_version" mapstructure:"account_version"`
	ConfirmationHeight         string `json:"confirmation_height" mapstructure:"confirmation_height"`
	ConfirmationHeightFrontier string `json:"confirmation_height_frontier" mapstructure:"confirmation_height_frontier"`
	Representative             string `json:"representative" mapstructure:"representative"`
	Weight                     string `json:"weight" mapstructure:"weight"`
	Pending                    string `json:"pending" mapstructure:"pending"`
	ConfirmedHeight            string `json:"confirmed_height" mapstructure:"confirmed_height"`
	ConfirmedFrontier          string `json:"confirmed_frontier" mapstructure:"confirmed_frontier"`
	Receivable                 string `json:"receivable" mapstructure:"receivable"`
	ConfirmedBalance           string `json:"confirmed_balance" mapstructure:"confirmed_balance"`
	ConfirmedPending           string `json:"confirmed_pending" mapstructure:"confirmed_pending"`
	ConfirmedReceivable        string `json:"confirmed_receivable" mapstructure:"confirmed_receivable"`
}
