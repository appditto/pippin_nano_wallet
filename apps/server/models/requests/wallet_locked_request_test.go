package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeWalletLockedRequest(t *testing.T) {
	encoded := `{"action":"wallet_locked","wallet":"1234"}`
	var decoded WalletLockedRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "wallet_locked", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)

	encoded = `{"action":"wallet_locked","wallet":""}`
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "wallet_locked", decoded.Action)
	assert.Equal(t, "", decoded.Wallet)
}

func TestMapStructureDecodeWalletLockedRequest(t *testing.T) {
	request := map[string]interface{}{
		"action": "wallet_locked",
		"wallet": "1234",
	}
	var decoded WalletLockedRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "wallet_locked", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
}
