package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeWalletCreateResponse(t *testing.T) {
	response := WalletCreateResponse{
		Wallet: "wallet",
	}
	encoded, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Equal(t, "{\"wallet\":\"wallet\"}", string(encoded))
}
