package responses

type AccountsFrontiersResponse struct {
	Frontiers *map[string]string `json:"frontiers,omitempty" mapstructure:"frontiers,omitempty"`
}
