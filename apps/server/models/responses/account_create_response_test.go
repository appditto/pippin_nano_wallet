package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeAccountCreateResponse(t *testing.T) {
	response := AccountCreateResponse{
		Account: "account",
	}
	encoded, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Equal(t, "{\"account\":\"account\"}", string(encoded))
}
