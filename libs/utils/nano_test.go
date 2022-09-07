package utils

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressToPub(t *testing.T) {
	address := "ban_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j"
	pub, err := AddressToPub(address)
	assert.Nil(t, err)
	assert.Equal(t, "dba12a8ed2702404483f5519c2b24e3c2f9cfc92310ab43255f3f6369e27a527", hex.EncodeToString(pub))
	address = "xrb_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j"
	pub, err = AddressToPub(address)
	assert.Nil(t, err)
	assert.Equal(t, "dba12a8ed2702404483f5519c2b24e3c2f9cfc92310ab43255f3f6369e27a527", hex.EncodeToString(pub))
	address = "nano_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j"
	pub, err = AddressToPub(address)
	assert.Nil(t, err)
	assert.Equal(t, "dba12a8ed2702404483f5519c2b24e3c2f9cfc92310ab43255f3f6369e27a527", hex.EncodeToString(pub))

	// Invalid
	address = "nano_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba511"
	pub, err = AddressToPub(address)
	assert.NotNil(t, err)
	assert.Equal(t, "", hex.EncodeToString(pub))
}

func TestGetAddresChecksum(t *testing.T) {
	address := "ban_3px37c9f6w361j65yoasrcs6wh3hmmyb6eacpis7dwzp8th4hbb9izgba51j"
	pub, err := AddressToPub(address)
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
