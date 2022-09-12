package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeWorkGenerateResponse(t *testing.T) {
	response := WorkGenerateResponse{
		Work: "1",
	}
	encoded, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Equal(t, "{\"work\":\"1\"}", string(encoded))
}
