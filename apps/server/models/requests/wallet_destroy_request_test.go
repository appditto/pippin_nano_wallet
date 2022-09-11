package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeWalletDestroyRequest(t *testing.T) {
	encoded := `{"action":"wallet_destroy","wallet":"1234"}`
	var decoded WalletDestroyRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "wallet_destroy", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
}

func TestMapStructureDecodeWalletDestroyRequest(t *testing.T) {
	request := map[string]interface{}{
		"action": "wallet_destroy",
		"wallet": "1234",
	}
	var decoded WalletDestroyRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "wallet_destroy", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
}
