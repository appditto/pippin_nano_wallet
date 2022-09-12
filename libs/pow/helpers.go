package pow

import (
	"encoding/binary"
	"encoding/hex"
	"strconv"

	"github.com/bbedward/nanopow"
	"golang.org/x/crypto/blake2b"
)

const (
	baseMaxUint64  = uint64(1<<64 - 1)
	baseDifficulty = baseMaxUint64 - uint64(0xfffffe0000000000)
)

// This is a helper to convert work multiplier to difficulty string representation
// BoomPoW takes a multiplier while the node/other work servers take the string
// Our base is banano or nano's receive, which would be 1x
// Nano's send would be 64x
func DifficultyFromMultiplier(multiplier int) uint64 {
	if multiplier < 1 {
		multiplier = 1
	}

	return baseMaxUint64 - (baseDifficulty / uint64(multiplier))
}

func DifficultyToString(difficulty uint64) string {
	return strconv.FormatUint(difficulty, 16)
}

func IsWorkValid(previous string, difficultyMultiplier int, w string) bool {
	difficult := DifficultyFromMultiplier(difficultyMultiplier)
	previousEnc, err := hex.DecodeString(previous)
	if err != nil {
		return false
	}
	wEnc, err := hex.DecodeString(w)
	if err != nil {
		return false
	}

	hash, err := blake2b.New(8, nil)
	if err != nil {
		return false
	}

	n := make([]byte, 8)
	copy(n, wEnc[:])

	reverse(n)
	hash.Write(n)
	hash.Write(previousEnc[:])

	return binary.LittleEndian.Uint64(hash.Sum(nil)) >= difficult
}

func reverse(v []byte) {
	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
		v[i], v[j] = v[j], v[i]
	}
}

func WorkToString(w nanopow.Work) string {
	n := make([]byte, 8)
	copy(n, w[:])

	return hex.EncodeToString(n)
}
