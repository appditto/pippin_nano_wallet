package responses

type WalletChangeSeedResponse struct {
	Success             string `json:"success" mapstructure:"success"`
	LastRestoredAccount string `json:"last_restored_account" mapstructure:"last_restored_account"`
	RestoredCount       int    `json:"restored_count" mapstructure:"restored_count"`
}
