package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeAccountInfoRequest(t *testing.T) {
	encoded := `{"action":"account_info","account":"abc","representative":true}`
	var decoded AccountInfoRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "account_info", decoded.Action)
	assert.Equal(t, "abc", decoded.Account)
	assert.True(t, *decoded.Representative)
	assert.Nil(t, decoded.Pending)
}

func TestMapStructureDecodeAccountInfoRequest(t *testing.T) {
	request := map[string]interface{}{
		"action":         "account_info",
		"account":        "abc",
		"representative": true,
	}
	var decoded AccountInfoRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "account_info", decoded.Action)
	assert.Equal(t, "abc", decoded.Account)
	assert.True(t, *decoded.Representative)
	assert.Nil(t, decoded.Pending)
}

// Test encoding
func TestEncodeAccountInfoRequest(t *testing.T) {
	encoded := `{"action":"account_info","account":"abc","representative":true}`
	rep := true
	req := AccountInfoRequest{
		AccountRequest: AccountRequest{
			BaseRequest: BaseRequest{Action: "account_info"},
			Account:     "abc",
		},
		Representative: &rep,
	}
	encodedActual, _ := json.Marshal(&req)
	assert.Equal(t, encoded, string(encodedActual))
}
