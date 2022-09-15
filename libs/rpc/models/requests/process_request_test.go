package requests

import (
	"encoding/json"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/utils"
	walletmodels "github.com/appditto/pippin_nano_wallet/libs/wallet/models"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestEncodeProcessRequest(t *testing.T) {
	stateBlock := walletmodels.StateBlock{
		Type:           "state",
		Hash:           "1234",
		Account:        "1",
		Previous:       "2",
		Representative: "3",
		Balance:        "4",
		Link:           "5",
	}
	request := ProcessRequest{
		BaseRequest: BaseRequest{
			Action: "process",
		},
		Block:     stateBlock,
		JsonBlock: true,
	}
	encoded, err := json.Marshal(request)
	assert.Nil(t, err)
	assert.Equal(t, `{"action":"process","json_block":true,"block":{"type":"state","hash":"1234","account":"1","previous":"2","representative":"3","balance":"4","link":"5","work":"","signature":""}}`, string(encoded))
}

func TestDecodeProcessRequest(t *testing.T) {
	encoded := `{"action":"process","json_block":"true","block":{"type":"state","hash":"1234","account":"1","previous":"2","representative":"3","balance":"4","link":"5","work":"","signature":""}}`
	var request ProcessRequest
	err := json.Unmarshal([]byte(encoded), &request)
	assert.Nil(t, err)
	assert.Equal(t, "process", request.Action)
	asBool, _ := utils.ToBool(request.JsonBlock)
	assert.True(t, asBool)
	assert.Equal(t, "state", request.Block.Type)
	assert.Equal(t, "1234", request.Block.Hash)
	assert.Equal(t, "1", request.Block.Account)
	assert.Equal(t, "2", request.Block.Previous)
	assert.Equal(t, "3", request.Block.Representative)
	assert.Equal(t, "4", request.Block.Balance)
	assert.Equal(t, "5", request.Block.Link)
	assert.Equal(t, "", request.Block.Work)
	assert.Equal(t, "", request.Block.Signature)
}

func TestMapStructureDecodeProcessRequest(t *testing.T) {
	request := map[string]interface{}{
		"action":     "process",
		"json_block": "true",
		"block": map[string]interface{}{
			"type":           "state",
			"hash":           "1234",
			"account":        "1",
			"previous":       "2",
			"representative": "3",
			"balance":        "4",
			"link":           "5",
			"work":           "",
			"signature":      "",
		},
	}
	var decoded ProcessRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "process", decoded.Action)
	asBool, _ := utils.ToBool(decoded.JsonBlock)
	assert.True(t, asBool)
	assert.Equal(t, "state", decoded.Block.Type)
	assert.Equal(t, "1234", decoded.Block.Hash)
	assert.Equal(t, "1", decoded.Block.Account)
	assert.Equal(t, "2", decoded.Block.Previous)
	assert.Equal(t, "3", decoded.Block.Representative)
	assert.Equal(t, "4", decoded.Block.Balance)
	assert.Equal(t, "5", decoded.Block.Link)
	assert.Equal(t, "", decoded.Block.Work)
	assert.Equal(t, "", decoded.Block.Signature)
}
