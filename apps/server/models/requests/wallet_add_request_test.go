package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeWalletAddRequest(t *testing.T) {
	encoded := `{"action":"wallet_add","key":"1234","wallet":"1234"}`
	var decoded WalletAddRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "wallet_add", decoded.Action)
	assert.Equal(t, "1234", decoded.Key)
	assert.Equal(t, "1234", decoded.Wallet)

	encoded = `{"action":"wallet_add","key":"","wallet":"12345"}`
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "wallet_add", decoded.Action)
	assert.Equal(t, "", decoded.Key)
	assert.Equal(t, "12345", decoded.Wallet)
}

func TestMapStructureDecodeWalletAddRequest(t *testing.T) {
	request := map[string]interface{}{
		"action": "wallet_add",
		"key":    "1234",
		"wallet": "1234",
	}
	var decoded WalletAddRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "wallet_add", decoded.Action)
	assert.Equal(t, "1234", decoded.Key)
	assert.Equal(t, "1234", decoded.Wallet)
}
