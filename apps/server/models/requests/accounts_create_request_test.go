package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeAccountsCreateRequest(t *testing.T) {
	encoded := `{"action":"accounts_create","wallet":"1234","count":10}`
	var decoded AccountsCreateRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "accounts_create", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
	assert.Equal(t, 10, decoded.Count)
}

func TestMapStructureDecodeAccountsCreateRequest(t *testing.T) {
	request := map[string]interface{}{
		"action": "accounts_create",
		"wallet": "1234",
		"count":  10,
	}
	var decoded AccountsCreateRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "accounts_create", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
	assert.Equal(t, 10, decoded.Count)
}
