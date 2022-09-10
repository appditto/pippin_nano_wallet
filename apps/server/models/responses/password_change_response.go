package responses

type PasswordChangeResponse struct {
	Changed string `json:"changed" mapstructure:"changed"`
}
