package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeAccountsPendingResponse(t *testing.T) {
	encoded := "{\n  \"blocks\" : {\n    \"nano_1111111111111111111111111111111111111111111111111117353trpda\": [\"142A538F36833D1CC78B94E11C766F75818F8B940771335C6C1B8AB880C5BB1D\"],\n    \"nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3\": [\"4C1FEEF0BEA7F50BE35489A1233FE002B212DEA554B55B1B470D78BD8F210C74\"]\n  }\n}"
	var decoded AccountsPendingResponse
	json.Unmarshal([]byte(encoded), &decoded)
	blocks := *decoded.Blocks
	assert.Len(t, blocks, 2)
	assert.Equal(t, "4C1FEEF0BEA7F50BE35489A1233FE002B212DEA554B55B1B470D78BD8F210C74", blocks["nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3"][0])
	assert.Equal(t, "142A538F36833D1CC78B94E11C766F75818F8B940771335C6C1B8AB880C5BB1D", blocks["nano_1111111111111111111111111111111111111111111111111117353trpda"][0])
}

func TestDecodeAccountsPendingResponseError(t *testing.T) {
	encoded := "{\"error\": \"Account not found\"}"
	var decoded AccountsPendingResponse
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Nil(t, decoded.Blocks)
}
