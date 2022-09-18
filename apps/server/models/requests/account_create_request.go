package requests

type AccountCreateRequest struct {
	BaseRequest `mapstructure:",squash"`
	Index       *interface{} `json:"index,omitempty" mapstructure:"index,omitempty"`
}
