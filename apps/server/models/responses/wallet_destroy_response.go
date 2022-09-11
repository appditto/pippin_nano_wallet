package responses

type WalletDestroyResponse struct {
	Destroyed string `json:"destroyed" mapstructure:"destroyed"`
}
