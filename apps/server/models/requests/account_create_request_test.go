package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeAccountCreateRequest(t *testing.T) {
	encoded := `{"action":"account_create","wallet":"1234","index":1}`
	var decoded AccountCreateRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "account_create", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
	assert.Equal(t, 1.0, *decoded.Index)
}

func TestMapStructureDecodeAccountCreateRequest(t *testing.T) {
	request := map[string]interface{}{
		"action": "account_create",
		"wallet": "1234",
		"index":  1,
	}
	var decoded AccountCreateRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "account_create", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
	assert.Equal(t, 1, *decoded.Index)
}
