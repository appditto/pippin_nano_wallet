package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeBaseRequest(t *testing.T) {
	encoded := `{"action":"account_create","wallet":"1234"}`
	var decoded BaseRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "account_create", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
}

func TestMapStructureDecodeBaseRequest(t *testing.T) {
	request := map[string]interface{}{
		"action": "account_create",
		"wallet": "1234",
	}
	var decoded BaseRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "account_create", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
}
