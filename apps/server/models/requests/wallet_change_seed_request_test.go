package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeWalletChangeSeedRequest(t *testing.T) {
	encoded := `{"action":"wallet_change_seed","seed":"sdasdas","wallet":"1234"}`
	var decoded WalletChangeSeedRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "wallet_change_seed", decoded.Action)
	assert.Equal(t, "sdasdas", decoded.Seed)
	assert.Equal(t, "1234", decoded.Wallet)
}

func TestMapStructureDecodeWalletChangeSeedRuest(t *testing.T) {
	request := map[string]interface{}{
		"action": "wallet_change_seed",
		"seed":   "sdasdas",
		"wallet": "1234",
	}
	var decoded WalletChangeSeedRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "wallet_change_seed", decoded.Action)
	assert.Equal(t, "sdasdas", decoded.Seed)
	assert.Equal(t, "1234", decoded.Wallet)
}
