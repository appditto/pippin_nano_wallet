package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeBaseRequest(t *testing.T) {
	encoded := "{\"action\":\"action\"}"
	var request BaseRequest
	err := json.Unmarshal([]byte(encoded), &request)
	assert.Nil(t, err)
	assert.Equal(t, "action", request.Action)
}

func TestMapStructureDecodeBaseRequest(t *testing.T) {
	request := map[string]interface{}{
		"action": "action",
	}
	var decoded BaseRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "action", decoded.Action)
}
