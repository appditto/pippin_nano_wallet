package models

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeWorkGenerateResponse(t *testing.T) {
	encoded := `{"work":"1234"}`
	var decoded WorkGenerateResponse
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "1234", decoded.Work)
}

func TestMapStructureDecodeWorkGenerateResponse(t *testing.T) {
	request := map[string]interface{}{
		"work": "def",
	}
	var decoded WorkGenerateResponse
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "def", decoded.Work)
}
