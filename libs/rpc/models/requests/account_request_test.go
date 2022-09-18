package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeAccounsBalancesRequest(t *testing.T) {
	encoded := `{"action":"account_balance","account":"abc"}`
	var decoded AccountRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "account_balance", decoded.Action)
	assert.Equal(t, "abc", decoded.Account)
}

func TestMapStructureDecodeAccountBalanceRequest(t *testing.T) {
	request := map[string]interface{}{
		"action":  "account_balance",
		"account": "abc",
	}
	var decoded AccountRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "account_balance", decoded.Action)
	assert.Equal(t, "abc", decoded.Account)
}

// Test encoding
func TestEncodeAccountBalancesRequest(t *testing.T) {
	encoded := `{"action":"account_balances","account":"abc"}`
	req := AccountRequest{
		BaseRequest: BaseRequest{Action: "account_balances"},
		Account:     "abc",
	}
	encodedActual, _ := json.Marshal(&req)
	assert.Equal(t, encoded, string(encodedActual))
}
