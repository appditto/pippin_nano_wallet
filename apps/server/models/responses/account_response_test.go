package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeAccountResponse(t *testing.T) {
	response := AccountResponse{
		Account: "account",
	}
	encoded, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Equal(t, "{\"account\":\"account\"}", string(encoded))
}
