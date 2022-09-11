package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test encoding
func TestEncodeWorkGenerateRequest(t *testing.T) {
	encoded := `{"action":"work_generate","hash":"abc","difficulty":"def"}`
	req := WorkGenerateRequest{
		WorkBaseRequest: WorkBaseRequest{
			Action: "work_generate",
			Hash:   "abc",
		},
		Difficulty: "def",
	}
	encodedActual, _ := json.Marshal(&req)
	assert.Equal(t, encoded, string(encodedActual))
}
