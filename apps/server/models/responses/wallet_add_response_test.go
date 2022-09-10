package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeWalletAddResponse(t *testing.T) {
	response := WalletAddResponse{
		Account: "acc",
	}
	encoded, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Equal(t, "{\"account\":\"acc\"}", string(encoded))
}
