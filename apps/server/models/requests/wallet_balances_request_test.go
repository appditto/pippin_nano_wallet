package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeWalletBalancesRequest(t *testing.T) {
	encoded := `{"action":"wallet_balances","wallet":"1234"}`
	var decoded WalletBalancesRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "wallet_balances", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
}

func TestMapStructureDecodeWalletBalancesRequest(t *testing.T) {
	request := map[string]interface{}{
		"action": "wallet_balances",
		"wallet": "1234",
	}
	var decoded WalletBalancesRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "wallet_balances", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
}
