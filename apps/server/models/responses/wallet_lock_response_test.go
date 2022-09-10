package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWalletLockResponse(t *testing.T) {
	response := WalletLockResponse{
		Locked: "1",
	}
	encoded, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Equal(t, "{\"locked\":\"1\"}", string(encoded))
}
