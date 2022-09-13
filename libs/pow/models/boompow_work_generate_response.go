package models

type BoompowWorkGenerateResponse struct {
	WorkGenerate string `json:"workGenerate" mapstructure:"workGenerate"`
}

type BoompowResponse struct {
	Data BoompowWorkGenerateResponse `json:"data" mapstructure:"data"`
}
