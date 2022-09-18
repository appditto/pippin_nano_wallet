package requests

type WalletRepresentativeSetRequest struct {
	BaseRequest            `mapstructure:",squash"`
	Representative         string       `json:"representative" mapstructure:"representative"`
	UpdateExistingAccounts *interface{} `json:"update_existing_accounts,omitempty" mapstructure:"update_existing_accounts,omitempty"`
}
