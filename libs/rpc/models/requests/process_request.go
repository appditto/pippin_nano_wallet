package requests

import walletmodels "github.com/appditto/pippin_nano_wallet/libs/wallet/models"

type ProcessRequest struct {
	BaseRequest `mapstructure:",squash"`
	// As an interface since the RPC spec allows strings or bools
	JsonBlock interface{}             `json:"json_block" mapstructure:"json_block"`
	Subtype   *string                 `json:"subtype,omitempty" mapstructure:"subtype,omitempty"`
	Block     walletmodels.StateBlock `json:"block" mapstructure:"block"`
}
