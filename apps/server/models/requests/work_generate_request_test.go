package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodeWorkGenerateRequest(t *testing.T) {
	encoded := `{"action":"work_generate","hash":"my hash","difficulty":"my difficulty","subtype":"my subtype","bpow_key":"my bpow key"}`
	var decoded WorkGenerateRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "work_generate", decoded.Action)
	assert.Equal(t, "my hash", decoded.Hash)
	assert.Equal(t, "my difficulty", decoded.Difficulty)
	assert.Equal(t, "my subtype", decoded.Subtype)
	assert.Nil(t, decoded.BlockAward)
	assert.Equal(t, "my bpow key", decoded.BpowKey)
}

func TestMapStructureDecodeWorkGenerateRequest(t *testing.T) {
	request := map[string]interface{}{
		"action":      "work_generate",
		"hash":        "my hash",
		"difficulty":  "my difficulty",
		"subtype":     "my subtype",
		"block_award": true,
		"bpow_key":    "my bpow key",
	}
	var decoded WorkGenerateRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "work_generate", decoded.Action)
	assert.Equal(t, "my hash", decoded.Hash)
	assert.Equal(t, "my difficulty", decoded.Difficulty)
	assert.Equal(t, "my subtype", decoded.Subtype)
	assert.Equal(t, true, *decoded.BlockAward)
	assert.Equal(t, "my bpow key", decoded.BpowKey)
}
