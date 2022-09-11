package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeAccountsBalancesRequest(t *testing.T) {
	encoded := `{"action":"accounts_balances","accounts":["abc","def"]}`
	var decoded AccountsRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "accounts_balances", decoded.Action)
	assert.Len(t, decoded.Accounts, 2)
	assert.Equal(t, "abc", decoded.Accounts[0])
	assert.Equal(t, "def", decoded.Accounts[1])
}

func TestMapStructureDecodeAccountBalancesRequest(t *testing.T) {
	request := map[string]interface{}{
		"action":   "accounts_balances",
		"accounts": []string{"abc", "def"},
	}
	var decoded AccountsRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "accounts_balances", decoded.Action)
	assert.Len(t, decoded.Accounts, 2)
	assert.Equal(t, "abc", decoded.Accounts[0])
	assert.Equal(t, "def", decoded.Accounts[1])
}

// Test encoding
func TestEncodeAccountsBalancesRequest(t *testing.T) {
	encoded := `{"action":"accounts_balances","accounts":["abc","def"]}`
	req := AccountsRequest{
		BaseRequest: BaseRequest{Action: "accounts_balances"},
		Accounts:    []string{"abc", "def"},
	}
	encodedActual, _ := json.Marshal(&req)
	assert.Equal(t, encoded, string(encodedActual))
}
