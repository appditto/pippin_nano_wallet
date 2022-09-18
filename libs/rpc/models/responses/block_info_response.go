package responses

import walletmodels "github.com/appditto/pippin_nano_wallet/libs/wallet/models"

type BlockInfoResponse struct {
	BlockAccount   string                  `json:"block_account" mapstructure:"block_account"`
	Amount         string                  `json:"amount" mapstructure:"amount"`
	Balance        string                  `json:"balance" mapstructure:"balance"`
	Height         string                  `json:"height" mapstructure:"height"`
	LocalTimestamp string                  `json:"local_timestamp" mapstructure:"local_timestamp"`
	Successor      string                  `json:"successor" mapstructure:"successor"`
	Confirmed      string                  `json:"confirmed" mapstructure:"confirmed"`
	Contents       walletmodels.StateBlock `json:"contents" mapstructure:"contents"`
	Subtype        string                  `json:"subtype" mapstructure:"subtype"`
}
