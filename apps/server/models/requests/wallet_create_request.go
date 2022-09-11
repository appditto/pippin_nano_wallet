package requests

type WalletCreateRequest struct {
	Action string  `json:"action" mapstructure:"action"`
	Seed   *string `json:"seed,omitempty" mapstructure:"seed,omitempty"`
}
