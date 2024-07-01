package requests

import (
	"encoding/json"
	"fmt"
)

type SendRequest struct {
	BaseRequest `mapstructure:",squash"`
	Source      string          `json:"source" mapstructure:"source"`
	Destination string          `json:"destination" mapstructure:"destination"`
	Amount      string          `json:"-"` // Exclude from automatic unmarshalling
	RawAmount   json.RawMessage `json:"amount" mapstructure:"amount"`
	ID          *string         `json:"id,omitempty" mapstructure:"id,omitempty"`
	Work        *string         `json:"work,omitempty" mapstructure:"work,omitempty"`
}

func (r *SendRequest) UnmarshalJSON(data []byte) error {
	type Alias SendRequest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var temp interface{}
	if err := json.Unmarshal(r.RawAmount, &temp); err != nil {
		return err
	}

	switch v := temp.(type) {
	case float64:
		r.Amount = fmt.Sprintf("%.0f", v) // Safely convert to string without scientific notation
	case int:
		r.Amount = fmt.Sprintf("%d", v)
	case string:
		r.Amount = v
	default:
		return fmt.Errorf("amount is of a type I don't know how to handle")
	}

	return nil
}
