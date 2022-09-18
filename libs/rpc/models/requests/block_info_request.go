package requests

type BlockInfoRequest struct {
	BaseRequest `mapstructure:",squash"`
	JsonBlock   bool   `json:"json_block" mapstructure:"json_block"`
	Hash        string `json:"hash" mapstructure:"hash"`
}
