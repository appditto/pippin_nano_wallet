package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeWalletContainsResponse(t *testing.T) {
	response := WalletContainsResponse{
		Exists: "0",
	}
	encoded, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Equal(t, "{\"exists\":\"0\"}", string(encoded))
}
