package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeAccountRepresentativeSetRequest(t *testing.T) {
	encoded := `{"action":"account_representative_set","wallet":"1234","account":"nano_1","representative":"nano_2"}`
	var decoded AccountRepresentativeSetRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "account_representative_set", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
	assert.Equal(t, "nano_1", decoded.Account)
	assert.Equal(t, "nano_2", decoded.Representative)
	assert.Nil(t, decoded.Work)
	assert.Nil(t, decoded.BpowKey)
}

func TestMapStructureDecodeAccountRepresentativeSetRequest(t *testing.T) {
	request := map[string]interface{}{
		"action":         "account_representative_set",
		"wallet":         "1234",
		"account":        "nano_1",
		"representative": "nano_2",
	}
	var decoded AccountRepresentativeSetRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "account_representative_set", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
	assert.Equal(t, "nano_1", decoded.Account)
	assert.Equal(t, "nano_2", decoded.Representative)
	assert.Nil(t, decoded.Work)
	assert.Nil(t, decoded.BpowKey)
}
