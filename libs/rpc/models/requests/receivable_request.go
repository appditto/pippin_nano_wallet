package requests

type ReceivableRequest struct {
	BaseRequest          `mapstructure:",squash"`
	Account              string `json:"account" mapstructure:"account"`
	Threshold            string `json:"threshold" mapstructure:"threshold"`
	IncludeOnlyConfirmed bool   `json:"include_only_confirmed" mapstructure:"include_only_confirmed"`
}
