package requests

type WalletChangeSeedRequest struct {
	BaseRequest `mapstructure:",squash"`
	Seed        string `json:"seed" mapstructure:"seed"`
}
