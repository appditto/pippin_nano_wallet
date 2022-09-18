package models

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/utils/ed25519"
	"github.com/mitchellh/mapstructure"

	"github.com/stretchr/testify/assert"
)

func TestComputeBlockHash(t *testing.T) {
	sb := StateBlock{
		Account:        "xrb_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j",
		Previous:       "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		Representative: "xrb_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j",
		Balance:        "1000000000000000000000000000000",
		Link:           "d9dd06646f96474a46c57c13677812305120be228f39964e222c06ab89f63745",
		Banano:         false,
	}

	// Hash
	err := sb.computeHash()
	assert.Nil(t, err)
	assert.Equal(t, "8ebeb9534a14e0b17b3cd4639721387dedac80789278b540ddbde2a0b267b6d0", sb.Hash)
}

func TestSignBlock(t *testing.T) {
	sb := StateBlock{
		Account:        "xrb_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j",
		Previous:       "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		Representative: "xrb_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j",
		Balance:        "1000000000000000000000000000000",
		Link:           "d9dd06646f96474a46c57c13677812305120be228f39964e222c06ab89f63745",
		Banano:         false,
	}

	// private key
	_, priv, err := ed25519.GenerateKey(strings.NewReader("9f729340e07eee69abac049c2fdd4a3c4b50e4672a2fabdf1ae295f2b4f3040b"))
	assert.Nil(t, err)
	err = sb.Sign(priv)
	assert.Nil(t, err)
	assert.Equal(t, "8ebeb9534a14e0b17b3cd4639721387dedac80789278b540ddbde2a0b267b6d0", sb.Hash)
	assert.Equal(t, "b580fa76c0b763aa8a8a90af8592155c9478554ce04c87b5fb115baae624eafa116e04fffb273405c0ffcff6dfb021526292ac4418f3988d7684e15e486f1409", sb.Signature)
}

func TestComputeBlockHashBanano(t *testing.T) {
	sb := StateBlock{
		Account:        "ban_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j",
		Previous:       "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		Representative: "ban_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j",
		Balance:        "1000000000000000000000000000000",
		Link:           "d9dd06646f96474a46c57c13677812305120be228f39964e222c06ab89f63745",
		Banano:         true,
	}

	// Hash
	err := sb.computeHash()
	assert.Nil(t, err)
	assert.Equal(t, "8ebeb9534a14e0b17b3cd4639721387dedac80789278b540ddbde2a0b267b6d0", sb.Hash)
}

func TestSignBlockBanano(t *testing.T) {
	sb := StateBlock{
		Account:        "ban_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j",
		Previous:       "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		Representative: "ban_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j",
		Balance:        "1000000000000000000000000000000",
		Link:           "d9dd06646f96474a46c57c13677812305120be228f39964e222c06ab89f63745",
		Banano:         true,
	}

	// private key
	_, priv, err := ed25519.GenerateKey(strings.NewReader("9f729340e07eee69abac049c2fdd4a3c4b50e4672a2fabdf1ae295f2b4f3040b"))
	assert.Nil(t, err)
	err = sb.Sign(priv)
	assert.Nil(t, err)
	assert.Equal(t, "8ebeb9534a14e0b17b3cd4639721387dedac80789278b540ddbde2a0b267b6d0", sb.Hash)
	assert.Equal(t, "b580fa76c0b763aa8a8a90af8592155c9478554ce04c87b5fb115baae624eafa116e04fffb273405c0ffcff6dfb021526292ac4418f3988d7684e15e486f1409", sb.Signature)
}

func TestMapstructureDecodeStateBlock(t *testing.T) {
	request := map[string]interface{}{
		"type":           "state",
		"hash":           "1234",
		"account":        "1",
		"previous":       "2",
		"representative": "3",
		"balance":        "4",
		"link":           "5",
	}
	var decoded StateBlock
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "state", decoded.Type)
	assert.Equal(t, "1234", decoded.Hash)
	assert.Equal(t, "1", decoded.Account)
	assert.Equal(t, "2", decoded.Previous)
	assert.Equal(t, "3", decoded.Representative)
	assert.Equal(t, "4", decoded.Balance)
	assert.Equal(t, "5", decoded.Link)
	assert.Equal(t, "", decoded.Signature)
	assert.Equal(t, "", decoded.Work)
}

func TestJsonDecodeStateBlock(t *testing.T) {
	encoded := "{\n\t\t\"type\":           \"state\",\n\t\t\"hash\":           \"1234\",\n\t\t\"account\":        \"1\",\n\t\t\"previous\":       \"2\",\n\t\t\"representative\": \"3\",\n\t\t\"balance\":        \"4\",\n\t\t\"link\":           \"5\"\n\t}"
	var decoded StateBlock
	err := json.Unmarshal([]byte(encoded), &decoded)
	assert.Nil(t, err)
	assert.Equal(t, "state", decoded.Type)
	assert.Equal(t, "1234", decoded.Hash)
	assert.Equal(t, "1", decoded.Account)
	assert.Equal(t, "2", decoded.Previous)
	assert.Equal(t, "3", decoded.Representative)
	assert.Equal(t, "4", decoded.Balance)
	assert.Equal(t, "5", decoded.Link)
	assert.Equal(t, "", decoded.Signature)
	assert.Equal(t, "", decoded.Work)
}

func TestJsonEncodeStateBlock(t *testing.T) {
	stateBlock := StateBlock{
		Type:           "state",
		Hash:           "1234",
		Account:        "1",
		Previous:       "2",
		Representative: "3",
		Balance:        "4",
		Link:           "5",
	}
	encoded, err := json.Marshal(stateBlock)
	assert.Nil(t, err)
	assert.Equal(t, `{"type":"state","hash":"1234","account":"1","previous":"2","representative":"3","balance":"4","link":"5","work":"","signature":""}`, string(encoded))
}
