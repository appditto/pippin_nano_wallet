package utils

import (
	"bytes"
	cryptorand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"io"

	"github.com/appditto/pippin_nano_wallet/libs/utils/ed25519"

	"golang.org/x/crypto/blake2b"
)

// Generates a seed, if rand is nil will use crypto/rand to generate secure seeds
func GenerateSeed(rand io.Reader) (string, error) {
	if rand == nil {
		rand = cryptorand.Reader
	}
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Generate a keypair from a seed at specified index
func KeypairFromSeed(seed string, index uint32) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	hash, err := blake2b.New(32, nil)
	if err != nil {
		return nil, nil, err
	}

	seed_data, err := hex.DecodeString(seed)
	if err != nil {
		return nil, nil, err
	}

	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, index)

	hash.Write(seed_data)
	hash.Write(bs)

	seed_bytes := hash.Sum(nil)
	pub, priv, err := ed25519.GenerateKey(bytes.NewReader(seed_bytes))

	if err != nil {
		return nil, nil, err
	}

	return pub, priv, nil
}
