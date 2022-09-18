package requests

type AccountInfoRequest struct {
	AccountRequest   `mapstructure:",squash"`
	Representative   *bool `json:"representative,omitempty" mapstructure:"representative,omitempty"`
	Pending          *bool `json:"pending,omitempty" mapstructure:"pending,omitempty"`
	Weight           *bool `json:"weight,omitempty" mapstructure:"weight,omitempty"`
	IncludeConfirmed *bool `json:"include_confirmed,omitempty" mapstructure:"include_confirmed,omitempty"`
}
