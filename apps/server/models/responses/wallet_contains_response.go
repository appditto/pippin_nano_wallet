package responses

type WalletContainsResponse struct {
	Exists string `json:"exists" mapstructure:"exists"`
}
