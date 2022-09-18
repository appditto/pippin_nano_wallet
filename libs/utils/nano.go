package utils

import (
	"encoding/base32"
	"errors"
	"fmt"

	"github.com/appditto/pippin_nano_wallet/libs/utils/ed25519"

	"golang.org/x/crypto/blake2b"
)

// nano uses a non-standard base32 character set.
const EncodeNano = "13456789abcdefghijkmnopqrstuwxyz"

var NanoEncoding = base32.NewEncoding(EncodeNano)

func AddressToPub(account string, banano bool) (public_key []byte, err error) {
	if len(account) < 64 {
		return nil, errors.New("Invalid account length")
	}
	address := string(account)

	if (!banano && address[:4] == "xrb_") || (banano && address[:4] == "ban_") {
		address = address[4:]
	} else if !banano && address[:5] == "nano_" {
		address = address[5:]
	} else {
		return nil, errors.New("Invalid address format")
	}
	// A valid nano address is 64 bytes long
	// First 5 are simply a hard-coded string nano_ for ease of use
	// The following 52 characters form the address, and the final
	// 8 are a checksum.
	// They are base 32 encoded with a custom encoding.
	if len(address) == 60 {
		// The nano address string is 260bits which doesn't fall on a
		// byte boundary. pad with zeros to 280bits.
		// (zeros are encoded as 1 in nano's 32bit alphabet)
		key_b32nano := "1111" + address[0:52]
		input_checksum := address[52:]

		key_bytes, err := NanoEncoding.DecodeString(key_b32nano)
		if err != nil {
			return nil, err
		}
		// strip off upper 24 bits (3 bytes). 20 padding was added by us,
		// 4 is unused as account is 256 bits.
		key_bytes = key_bytes[3:]

		// nano checksum is calculated by hashing the key and reversing the bytes
		valid := NanoEncoding.EncodeToString(GetAddressChecksum(key_bytes)) == input_checksum
		if valid {
			return key_bytes, nil
		} else {
			return nil, errors.New("Invalid address checksum")
		}
	}

	return nil, errors.New("Invalid address format")
}

func PubKeyToAddress(pub ed25519.PublicKey, banano bool) string {
	// Pubkey is 256bits, base32 must be multiple of 5 bits
	// to encode properly.
	// Pad the start with 0's and strip them off after base32 encoding
	padded := append([]byte{0, 0, 0}, pub...)
	address := NanoEncoding.EncodeToString(padded)[4:]
	checksum := NanoEncoding.EncodeToString(GetAddressChecksum(pub))

	var prefix string
	if banano {
		prefix = "ban_"
	} else {
		prefix = "nano_"
	}

	return fmt.Sprintf("%s%s%s", prefix, address, checksum)
}

func GetAddressChecksum(pub ed25519.PublicKey) []byte {
	hash, err := blake2b.New(5, nil)
	if err != nil {
		panic("Unable to create hash")
	}

	hash.Write(pub)
	return Reversed(hash.Sum(nil))
}

func Reversed(str []byte) (result []byte) {
	for i := len(str) - 1; i >= 0; i-- {
		result = append(result, str[i])
	}
	return result
}
