package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeReceiveAllResponse(t *testing.T) {
	response := ReceiveAllResponse{
		Received: 10,
	}
	encoded, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Equal(t, "{\"received\":10}", string(encoded))
}
