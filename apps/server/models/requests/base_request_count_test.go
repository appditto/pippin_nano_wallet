package requests

import (
	"encoding/json"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeBaseWithCountRequest(t *testing.T) {
	encoded := `{"action":"account_list","wallet":"1234"}`
	var decoded BaseRequestWithCount
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "account_list", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
	assert.Nil(t, decoded.Count)

	encoded = `{"action":"account_list","wallet":"1234","count":"10"}`
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "account_list", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
	count, _ := utils.ToInt(*decoded.Count)
	assert.Equal(t, 10, count)
}

func TestMapStructureDecodeWalletBaseWithCountRequest(t *testing.T) {
	request := map[string]interface{}{
		"action": "account_list",
		"wallet": "1234",
		"count":  "1",
	}
	var decoded BaseRequestWithCount
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "account_list", decoded.Action)
	assert.Equal(t, "1234", decoded.Wallet)
	count, _ := utils.ToInt(*decoded.Count)
	assert.Equal(t, 1, count)
}
