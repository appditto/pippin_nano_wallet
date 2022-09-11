package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeAccountsFrontiersResponse(t *testing.T) {
	encoded := "{\n  \"frontiers\" : {\n    \"nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3\": \"791AF413173EEE674A6FCF633B5DFC0F3C33F397F0DA08E987D9E0741D40D81A\",\n    \"nano_3i1aq1cchnmbn9x5rsbap8b15akfh7wj7pwskuzi7ahz8oq6cobd99d4r3b7\": \"6A32397F4E95AF025DE29D9BF1ACE864D5404362258E06489FABDBA9DCCC046F\"\n  }\n}"
	var decoded AccountsFrontiersResponse
	json.Unmarshal([]byte(encoded), &decoded)
	frontiers := *decoded.Frontiers
	assert.Equal(t, "791AF413173EEE674A6FCF633B5DFC0F3C33F397F0DA08E987D9E0741D40D81A", frontiers["nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3"])
	assert.Equal(t, "6A32397F4E95AF025DE29D9BF1ACE864D5404362258E06489FABDBA9DCCC046F", frontiers["nano_3i1aq1cchnmbn9x5rsbap8b15akfh7wj7pwskuzi7ahz8oq6cobd99d4r3b7"])
}

func TestDecodeAccountsFrontiersResponseError(t *testing.T) {
	encoded := "{\"error\": \"Account not found\"}"
	var decoded AccountsFrontiersResponse
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Nil(t, decoded.Frontiers)
}
