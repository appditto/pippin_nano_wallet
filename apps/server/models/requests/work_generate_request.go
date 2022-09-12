package requests

type WorkGenerateRequest struct {
	Action     string `json:"action" mapstructure:"action"`
	Hash       string `json:"hash" mapstructure:"hash"`
	Difficulty string `json:"difficulty" mapstructure:"difficulty"`
	Subtype    string `json:"subtype" mapstructure:"subtype"`
	BlockAward *bool  `json:"block_award,omitempty" mapstructure:"block_award,omitempty"`
	BpowKey    string `json:"bpow_key" mapstructure:"bpow_key"`
}
