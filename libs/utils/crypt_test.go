package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var aesCrypt *AESCrypt
var secondAesCrypt *AESCrypt

func init() {
	aesCrypt = NewAesCrypt("mypassword")
	secondAesCrypt = NewAesCrypt("mysecondpassword")
}

func TestEncryptMessage(t *testing.T) {
	encrypted, err := aesCrypt.Encrypt("my message")
	assert.Nil(t, err)
	assert.NotEqual(t, "my message", encrypted)
}

func TestDecodeMessage(t *testing.T) {
	encrypted, _ := aesCrypt.Encrypt("my message")
	decrypted, err := aesCrypt.Decrypt(encrypted)
	assert.Nil(t, err)
	assert.Equal(t, "my message", decrypted)
}

func TestDecodeMessageNotEqualDifferentKeys(t *testing.T) {
	encrypted, _ := aesCrypt.Encrypt("my message")
	encryptedTwo, _ := secondAesCrypt.Encrypt("my message")
	assert.NotEqual(t, encrypted, encryptedTwo)
	decrypted, err := aesCrypt.Decrypt(encryptedTwo)
	assert.NotNil(t, err)
	assert.Equal(t, "", decrypted)
	decryptedTwo, err := secondAesCrypt.Decrypt(encrypted)
	assert.NotNil(t, err)
	assert.Equal(t, "", decryptedTwo)
}
