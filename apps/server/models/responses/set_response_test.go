package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeSetResponse(t *testing.T) {
	response := SetResponse{
		Set: "1",
	}
	encoded, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Equal(t, "{\"set\":\"1\"}", string(encoded))
}
