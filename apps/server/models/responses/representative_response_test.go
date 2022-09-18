package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeRepresentativeResponse(t *testing.T) {
	response := RepresentativeResponse{
		Representative: "1",
	}
	encoded, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Equal(t, "{\"representative\":\"1\"}", string(encoded))
}
