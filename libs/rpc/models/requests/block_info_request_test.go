package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestEncodeBlockInfoRequest(t *testing.T) {
	request := BlockInfoRequest{
		BaseRequest: BaseRequest{
			Action: "block_info",
		},
		JsonBlock: true,
		Hash:      "1234",
	}
	encoded, err := json.Marshal(request)
	assert.Nil(t, err)
	assert.Equal(t, "{\"action\":\"block_info\",\"json_block\":true,\"hash\":\"1234\"}", string(encoded))
}

func TestDecodeBlockInfoRequest(t *testing.T) {
	encoded := "{\"action\":\"block_info\",\"json_block\":true,\"hash\":\"1234\"}"
	var request BlockInfoRequest
	err := json.Unmarshal([]byte(encoded), &request)
	assert.Nil(t, err)
	assert.Equal(t, "block_info", request.Action)
	assert.Equal(t, "1234", request.Hash)
	assert.True(t, request.JsonBlock)
}

func TestMapStructureDecodeBlockInfoRequest(t *testing.T) {
	request := map[string]interface{}{
		"action":     "block_info",
		"json_block": true,
		"hash":       "1234",
	}
	var decoded BlockInfoRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "block_info", decoded.Action)
	assert.Equal(t, "1234", decoded.Hash)
	assert.True(t, decoded.JsonBlock)
}
