package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWalletDestroyResponse(t *testing.T) {
	response := WalletDestroyResponse{
		Destroyed: "1",
	}
	encoded, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Equal(t, "{\"destroyed\":\"1\"}", string(encoded))
}
