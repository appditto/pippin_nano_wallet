package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeWalletContainsRequest(t *testing.T) {
	encoded := `{"action":"wallet_contains","wallet":"1234","account":"5555"}`
	var decoded WalletContainsRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "wallet_contains", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
	assert.Equal(t, "5555", decoded.Account)
}

func TestMapStructureDecodeWalletContainsRequest(t *testing.T) {
	request := map[string]interface{}{
		"action":  "wallet_contains",
		"wallet":  "1234",
		"account": "5555",
	}
	var decoded WalletContainsRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "wallet_contains", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
	assert.Equal(t, "5555", decoded.Account)
}
