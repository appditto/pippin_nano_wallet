package responses

type WalletLockResponse struct {
	Locked string `json:"locked" mapstructure:"locked"`
}
