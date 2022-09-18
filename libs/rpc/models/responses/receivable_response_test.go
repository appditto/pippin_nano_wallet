package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecoddReceivableResponse(t *testing.T) {
	encoded := "{\n  \"blocks\" : {\n    \"000D1BAEC8EC208142C99059B393051BAC8380F9B5A2E6B2489A277D81789F3F\": \"6000000000000000000000000000000\"\n  }\n}"

	var decoded ReceivableResponse
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "6000000000000000000000000000000", decoded.Blocks["000D1BAEC8EC208142C99059B393051BAC8380F9B5A2E6B2489A277D81789F3F"])
}
