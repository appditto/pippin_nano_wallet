package responses

type WalletLockedResponse struct {
	Locked string `json:"locked" mapstructure:"locked"`
}
