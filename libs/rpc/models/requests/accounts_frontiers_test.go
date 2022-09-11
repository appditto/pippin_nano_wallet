package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeAccountsFrontiersRequest(t *testing.T) {
	encoded := `{"action":"accounts_frontiers","accounts":["abc","def"]}`
	var decoded AccountsFrontiersRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "accounts_frontiers", decoded.Action)
	assert.Len(t, decoded.Accounts, 2)
	assert.Equal(t, "abc", decoded.Accounts[0])
	assert.Equal(t, "def", decoded.Accounts[1])
}

func TestMapStructureDecodeAccountFrontiersRequest(t *testing.T) {
	request := map[string]interface{}{
		"action":   "accounts_frontiers",
		"accounts": []string{"abc", "def"},
	}
	var decoded AccountsFrontiersRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "accounts_frontiers", decoded.Action)
	assert.Len(t, decoded.Accounts, 2)
	assert.Equal(t, "abc", decoded.Accounts[0])
	assert.Equal(t, "def", decoded.Accounts[1])
}

// Test encoding
func TestEncodeAccountsFrontiersRequest(t *testing.T) {
	encoded := `{"action":"accounts_frontiers","accounts":["abc","def"]}`
	req := AccountsFrontiersRequest{
		BaseRequest: BaseRequest{Action: "accounts_frontiers"},
		Accounts:    []string{"abc", "def"},
	}
	encodedActual, _ := json.Marshal(&req)
	assert.Equal(t, encoded, string(encodedActual))
}
