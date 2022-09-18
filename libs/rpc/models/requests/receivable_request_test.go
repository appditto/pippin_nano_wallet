package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestEncodeReceivableRequest(t *testing.T) {
	request := ReceivableRequest{
		BaseRequest: BaseRequest{
			Action: "receivable",
		},
		Account:              "abcd",
		Threshold:            "1234",
		IncludeOnlyConfirmed: true,
	}
	encoded, err := json.Marshal(request)
	assert.Nil(t, err)
	assert.Equal(t, "{\"action\":\"receivable\",\"account\":\"abcd\",\"threshold\":\"1234\",\"include_only_confirmed\":true}", string(encoded))
}

func TestDecodeReceivableRequest(t *testing.T) {
	encoded := "{\"action\":\"receivable\",\"account\":\"abcd\",\"threshold\":\"1234\",\"include_only_confirmed\":true}"
	var request ReceivableRequest
	err := json.Unmarshal([]byte(encoded), &request)
	assert.Nil(t, err)
	assert.Equal(t, "receivable", request.Action)
	assert.Equal(t, "abcd", request.Account)
	assert.Equal(t, "1234", request.Threshold)
	assert.True(t, request.IncludeOnlyConfirmed)
}

func TestMapStructureDecodeReceivableRequest(t *testing.T) {
	request := map[string]interface{}{
		"action":                 "receivable",
		"account":                "abcd",
		"threshold":              "1234",
		"include_only_confirmed": true,
	}
	var decoded ReceivableRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "receivable", decoded.Action)
	assert.Equal(t, "abcd", decoded.Account)
	assert.Equal(t, "1234", decoded.Threshold)
	assert.True(t, decoded.IncludeOnlyConfirmed)
}
