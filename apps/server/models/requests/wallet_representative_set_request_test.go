package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeWalletRepresentativeSetRequest(t *testing.T) {
	encoded := `{"action":"wallet_representative_set","wallet":"1234","representative":"nano_2"}`
	var decoded WalletRepresentativeSetRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "wallet_representative_set", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
	assert.Equal(t, "nano_2", decoded.Representative)
	assert.Nil(t, decoded.BpowKey)
}

func TestMapStructureDecodeWalletRepresentativeSetRequest(t *testing.T) {
	request := map[string]interface{}{
		"action":         "wallet_representative_set",
		"wallet":         "1234",
		"representative": "nano_2",
	}
	var decoded WalletRepresentativeSetRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "wallet_representative_set", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
	assert.Equal(t, "nano_2", decoded.Representative)
	assert.Nil(t, decoded.BpowKey)
}
