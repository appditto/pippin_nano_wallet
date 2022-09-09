package utils

import "encoding/hex"

func Validate64HexHash(hash string) bool {
	if len(hash) != 64 {
		return false
	}
	_, err := hex.DecodeString(hash)
	if err != nil {
		return false
	}
	return true
}
