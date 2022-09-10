package requests

type AccountsCreateRequest struct {
	AccountCreateRequest `mapstructure:",squash"`
	// Make it an interface so it can be an int or a string
	Count interface{} `json:"count" mapstructure:"count"`
}
