package requests

import (
	"encoding/json"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeAccountsCreateRequest(t *testing.T) {
	encoded := `{"action":"accounts_create","wallet":"1234","count":10}`
	var decoded AccountsCreateRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "accounts_create", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
	count, _ := utils.ToInt(decoded.Count)
	assert.Equal(t, 10, count)

	encoded = `{"action":"accounts_create","wallet":"1234","count":"10"}`
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "accounts_create", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
	count, err := utils.ToInt(decoded.Count)
	assert.Nil(t, err)
	assert.Equal(t, 10, count)
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
	count, _ := utils.ToInt(decoded.Count)
	assert.Equal(t, 10, count)
}
