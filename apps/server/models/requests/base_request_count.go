package requests

// This is a common signature for most requests, the base + count
type BaseRequestWithCount struct {
	BaseRequest `mapstructure:",squash"`
	Count       *interface{} `json:"count,omitempty" mapstructure:"count,omitempty"`
}
