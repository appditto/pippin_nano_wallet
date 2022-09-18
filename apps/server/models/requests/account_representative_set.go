package requests

type AccountRepresentativeSetRequest struct {
	BaseRequest    `mapstructure:",squash"`
	Account        string  `json:"account" mapstructure:"account"`
	Representative string  `json:"representative" mapstructure:"representative"`
	Work           *string `json:"work,omitempty" mapstructure:"work,omitempty"`
}
