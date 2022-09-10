package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeWalletLockRequest(t *testing.T) {
	encoded := `{"action":"wallet_lock","wallet":"1234"}`
	var decoded WalletLockRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "wallet_lock", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)

	encoded = `{"action":"wallet_lock","wallet":""}`
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "wallet_lock", decoded.Action)
	assert.Equal(t, "", decoded.Wallet)
}

func TestMapStructureDecodeWalletLockRequest(t *testing.T) {
	request := map[string]interface{}{
		"action": "wallet_lock",
		"wallet": "1234",
	}
	var decoded WalletLockRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "wallet_lock", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
}
