package models

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/utils/ed25519"

	"github.com/stretchr/testify/assert"
)

func TestComputeBlockHash(t *testing.T) {
	previous, _ := hex.DecodeString("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	link, _ := hex.DecodeString("d9dd06646f96474a46c57c13677812305120be228f39964e222c06ab89f63745")
	bal, _ := big.NewInt(0).SetString("1000000000000000000000000000000", 10)

	sb := StateBlock{
		Account:        "xrb_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j",
		Previous:       previous,
		Representative: "xrb_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j",
		Balance:        *bal,
		Link:           link,
	}

	// Hash
	err := sb.computeHash()
	assert.Nil(t, err)
	assert.Equal(t, "8ebeb9534a14e0b17b3cd4639721387dedac80789278b540ddbde2a0b267b6d0", hex.EncodeToString(sb.Hash))
}

func TestSignBlock(t *testing.T) {
	previous, _ := hex.DecodeString("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	link, _ := hex.DecodeString("d9dd06646f96474a46c57c13677812305120be228f39964e222c06ab89f63745")
	bal, _ := big.NewInt(0).SetString("1000000000000000000000000000000", 10)

	sb := StateBlock{
		Account:        "xrb_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j",
		Previous:       previous,
		Representative: "xrb_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j",
		Balance:        *bal,
		Link:           link,
	}

	// private key
	_, priv, err := ed25519.GenerateKey(strings.NewReader("9f729340e07eee69abac049c2fdd4a3c4b50e4672a2fabdf1ae295f2b4f3040b"))
	assert.Nil(t, err)
	err = sb.Sign(priv)
	assert.Nil(t, err)
	assert.Equal(t, "8ebeb9534a14e0b17b3cd4639721387dedac80789278b540ddbde2a0b267b6d0", hex.EncodeToString(sb.Hash))
	assert.Equal(t, "b580fa76c0b763aa8a8a90af8592155c9478554ce04c87b5fb115baae624eafa116e04fffb273405c0ffcff6dfb021526292ac4418f3988d7684e15e486f1409", hex.EncodeToString(sb.Signature))
}
