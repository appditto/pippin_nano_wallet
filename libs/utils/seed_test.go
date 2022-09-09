package utils

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSeed(t *testing.T) {
	seed, _ := GenerateSeed(nil)
	assert.Len(t, seed, 64)

	// Test predictable seed
	seed, _ = GenerateSeed(strings.NewReader("9f729340e07eee69abac049c2fdd4a3c4b50e4672a2fabdf1ae295f2b4f3040b"))
	assert.Len(t, seed, 64)
	assert.Equal(t, "3966373239333430653037656565363961626163303439633266646434613363", seed)
}

func TestKeypairFromSeed(t *testing.T) {
	seed := "E11A48D701EA1F8A66A4EB587CDC8808D726FE75B325DF204F62CA2B43F9ADA1"
	pub, priv, err := KeypairFromSeed(seed, 0)
	assert.Nil(t, err)
	// Hex representation
	privHex := hex.EncodeToString(priv)
	// First 32 bytes are the private key
	assert.Equal(t, "67EDBC8F904091738DF33B4B6917261DB91DD9002D3985A7BA090345264A46C6", strings.ToUpper(privHex[:64]))
	// Last 32 bytes are the public key
	pubHex := hex.EncodeToString(pub)
	assert.Equal(t, strings.ToUpper(privHex[64:]), strings.ToUpper(pubHex))
}
