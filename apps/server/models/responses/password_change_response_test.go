package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodePasswordChangeResponse(t *testing.T) {
	response := PasswordChangeResponse{
		Changed: "1",
	}
	encoded, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Equal(t, "{\"changed\":\"1\"}", string(encoded))
}
