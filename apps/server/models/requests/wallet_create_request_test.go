package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeWalletCreateRequest(t *testing.T) {
	encoded := `{"action":"wallet_create"}`
	var decoded WalletCreateRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "wallet_create", decoded.Action)

	encoded = `{"action":"wallet_create", "seed":"my seed"}`
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "wallet_create", decoded.Action)
	assert.Equal(t, "my seed", *decoded.Seed)
}

func TestMapStructureDecodeWalletCreateRequest(t *testing.T) {
	request := map[string]interface{}{
		"action": "wallet_create",
		"seed":   "my seed",
	}
	var decoded WalletCreateRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "wallet_create", decoded.Action)
	assert.Equal(t, "my seed", *decoded.Seed)

	var decodedNoSeed WalletCreateRequest
	request = map[string]interface{}{
		"action": "wallet_create",
	}
	mapstructure.Decode(request, &decodedNoSeed)
	assert.Equal(t, "wallet_create", decodedNoSeed.Action)
	assert.Nil(t, decodedNoSeed.Seed)
}
