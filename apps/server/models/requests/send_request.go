package requests

import (
	"bytes"
	"encoding/json"
)

type SendRequest struct {
	BaseRequest `mapstructure:",squash"`
	Source      string  `json:"source" mapstructure:"source"`
	Destination string  `json:"destination" mapstructure:"destination"`
	Amount      string  `json:"amount" mapstructure:"amount"`
	ID          *string `json:"id,omitempty" mapstructure:"id,omitempty"`
	Work        *string `json:"work,omitempty" mapstructure:"work,omitempty"`
}

func (r *SendRequest) UnmarshalJSON(data []byte) error {
	type Alias SendRequest
	aux := struct {
		*Alias
		Amount json.Number `json:"amount"`
	}{
		Alias: (*Alias)(r),
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	if err := decoder.Decode(&aux); err != nil {
		return err
	}

	r.Amount = aux.Amount.String()
	return nil
}
