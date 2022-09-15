package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeProcessResponse(t *testing.T) {
	encoded := "{\"hash\":\"1234\"}"

	var decoded ProcessResponse
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "1234", decoded.Hash)
}
