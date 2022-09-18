package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeBlockResponse(t *testing.T) {
	response := BlockResponse{
		Block: "1234",
	}
	encoded, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Equal(t, "{\"block\":\"1234\"}", string(encoded))
}
