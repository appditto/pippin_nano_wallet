package utils

import (
	"encoding/hex"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/utils/ed25519"

	"github.com/stretchr/testify/assert"
)

func TestAddressToPub(t *testing.T) {
	address := "ban_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j"
	pub, err := AddressToPub(address, true)
	assert.Nil(t, err)
	assert.Equal(t, "dba12a8ed2702404483f5519c2b24e3c2f9cfc92310ab43255f3f6369e27a527", hex.EncodeToString(pub))
	address = "xrb_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j"
	pub, err = AddressToPub(address, false)
	assert.Nil(t, err)
	assert.Equal(t, "dba12a8ed2702404483f5519c2b24e3c2f9cfc92310ab43255f3f6369e27a527", hex.EncodeToString(pub))
	address = "nano_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j"
	pub, err = AddressToPub(address, false)
	assert.Nil(t, err)
	assert.Equal(t, "dba12a8ed2702404483f5519c2b24e3c2f9cfc92310ab43255f3f6369e27a527", hex.EncodeToString(pub))

	// Invalid
	address = "nano_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba511"
	pub, err = AddressToPub(address, false)
	assert.NotNil(t, err)
	assert.Equal(t, "", hex.EncodeToString(pub))

	// Test banano vs nano modes
	address = "ban_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j"
	pub, err = AddressToPub(address, false)
	assert.NotNil(t, err)
	address = "xrb_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j"
	pub, err = AddressToPub(address, true)
	assert.NotNil(t, err)
	address = "nano_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j"
	pub, err = AddressToPub(address, true)
	assert.NotNil(t, err)
}

func TestPubkeyToAddress(t *testing.T) {
	pubkey := "58E3EC60070DD5D991B899E4BAB6CFD97657AB79A388A9276E4456108E13D6BB"
	pub, _ := hex.DecodeString(pubkey)
	pubEd25519 := ed25519.PublicKey(pub)

	address := PubKeyToAddress(pubEd25519, false)
	assert.Equal(t, "nano_1p95xji1g5gou8auj8h6qcuezpdpcyoqmawao6mpwj4p44939oouoturkggc", address)

	// Banano
	address = PubKeyToAddress(pubEd25519, true)
	assert.Equal(t, "ban_1p95xji1g5gou8auj8h6qcuezpdpcyoqmawao6mpwj4p44939oouoturkggc", address)
}

func TestGetAddresChecksum(t *testing.T) {
	address := "ban_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j"
	pub, err := AddressToPub(address, true)
	assert.Nil(t, err)
	checksum := GetAddressChecksum(pub)
	assert.Equal(t, "87dc940c11", hex.EncodeToString(checksum))
}

func TestReversed(t *testing.T) {
	arr := []byte{1, 2, 3}
	reversed := Reversed(arr)
	assert.Equal(t, reversed[0], uint8(3))
	assert.Equal(t, reversed[1], uint8(2))
	assert.Equal(t, reversed[2], uint8(1))
}
