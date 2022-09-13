package models

type BoompowWorkGenerateRequestVariables struct {
	Hash                 string `json:"hash" mapstructure:"hash"`
	DifficultyMultiplier int    `json:"difficultyMultiplier" mapstructure:"difficultyMultiplier"`
	BlockAward           bool   `json:"blockAward" mapstructure:"blockAward"`
}

type BoompowWorkGenerateRequest struct {
	Query     string                              `json:"query" mapstructure:"query"`
	Variables BoompowWorkGenerateRequestVariables `json:"variables" mapstructure:"variables"`
}
