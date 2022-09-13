package models

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeBoompowWorkGenerateResponse(t *testing.T) {
	encoded := `{"data":{"workGenerate":"1234"}}`
	var decoded BoompowResponse
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "1234", decoded.Data.WorkGenerate)
}

func TestMapDestructureBoompowWorkGenerateResponse(t *testing.T) {
	request := map[string]interface{}{
		"data": map[string]interface{}{
			"workGenerate": "def",
		},
	}
	var decoded BoompowResponse
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "def", decoded.Data.WorkGenerate)
}
