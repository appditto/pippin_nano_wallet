package models

type WorkGenerateRequest struct {
	WorkBaseRequest
	Difficulty string `json:"difficulty" mapstructure:"difficulty"`
}
