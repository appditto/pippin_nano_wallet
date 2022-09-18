package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeReceiveRequest(t *testing.T) {
	encoded := `{"action":"receive","wallet":"1234","account":"nano_1","bpow_key":"abc"}`
	var decoded ReceiveRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "receive", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
	assert.Equal(t, "nano_1", decoded.Account)
	assert.Equal(t, "abc", *decoded.BpowKey)
	assert.Nil(t, decoded.Work)
}

func TestMapStructureDecodeReceiveRequest(t *testing.T) {
	request := map[string]interface{}{
		"action":   "receive",
		"wallet":   "1234",
		"account":  "nano_1",
		"bpow_key": "abc",
	}
	var decoded ReceiveRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "receive", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
	assert.Equal(t, "nano_1", decoded.Account)
	assert.Equal(t, "abc", *decoded.BpowKey)
	assert.Nil(t, decoded.Work)
}
