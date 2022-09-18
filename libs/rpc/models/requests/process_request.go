package requests

import walletmodels "github.com/appditto/pippin_nano_wallet/libs/wallet/models"

type ProcessRequest struct {
	BaseRequest `mapstructure:",squash"`
	JsonBlock   bool                    `json:"json_block" mapstructure:"json_block"`
	Subtype     *string                 `json:"subtype,omitempty" mapstructure:"subtype,omitempty"`
	Block       walletmodels.StateBlock `json:"block" mapstructure:"block"`
}
