package models

type WorkBaseRequest struct {
	Action string `json:"action" mapstructure:"action"`
	Hash   string `json:"hash" mapstructure:"hash"`
}
