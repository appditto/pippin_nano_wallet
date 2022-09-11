package requests

type BaseRequest struct {
	Action string `json:"action" mapstructure:"action"`
}
