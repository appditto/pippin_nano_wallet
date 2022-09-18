package requests

type SendRequest struct {
	BaseRequest `mapstructure:",squash"`
	Source      string  `json:"source" mapstructure:"source"`
	Destination string  `json:"destination" mapstructure:"destination"`
	Amount      string  `json:"amount" mapstructure:"amount"`
	ID          *string `json:"id,omitempty" mapstructure:"id,omitempty"`
	Work        *string `json:"work,omitempty" mapstructure:"work,omitempty"`
}
