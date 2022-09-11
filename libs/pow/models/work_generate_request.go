package models

type WorkGenerateRequest struct {
	Action     string `json:"action" mapstructure:"action"`
	Hash       string `json:"hash" mapstructure:"hash"`
	Difficulty string `json:"difficulty" mapstructure:"difficulty"`
}
