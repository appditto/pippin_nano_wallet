package requests

type WalletCreateRequest struct {
	BaseRequest `mapstructure:",squash"`
	Seed        *string `json:"seed,omitempty" mapstructure:"seed,omitempty"`
}
