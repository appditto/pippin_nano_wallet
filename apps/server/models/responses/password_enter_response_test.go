package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodePasswordEnterResponse(t *testing.T) {
	response := PasswordEnterResponse{
		Valid: "0",
	}
	encoded, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Equal(t, "{\"valid\":\"0\"}", string(encoded))
}
