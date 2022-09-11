package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeWorkBaseRequest(t *testing.T) {
	encoded := `{"action":"work_cancel","hash":"abc"}`
	req := WorkBaseRequest{
		Action: "work_cancel",
		Hash:   "abc",
	}
	encodedActual, _ := json.Marshal(&req)
	assert.Equal(t, encoded, string(encodedActual))
}
