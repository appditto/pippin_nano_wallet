package requests

type AccountsFrontiersRequest struct {
	BaseRequest `mapstructure:",squash"`
	Accounts    []string `json:"accounts" mapstructure:"accounts"`
}
