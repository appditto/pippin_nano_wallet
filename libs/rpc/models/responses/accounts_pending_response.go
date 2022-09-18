package responses

type AccountsPendingResponse struct {
	Blocks *map[string][]string `json:"blocks" mapstructure:"blocks"`
}
