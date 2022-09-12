package pow

import (
	"testing"

	"github.com/bbedward/nanopow"
	"github.com/stretchr/testify/assert"
)

func TestDifficultyFromMultiplier(t *testing.T) {
	assert.Equal(t, uint64(0xfffffe0000000000), DifficultyFromMultiplier(0))
	assert.Equal(t, uint64(0xfffffe0000000000), DifficultyFromMultiplier(1))
	assert.Equal(t, uint64(0xfffffff800000000), DifficultyFromMultiplier(64))
}

func TestDifficultToString(t *testing.T) {
	assert.Equal(t, "fffffe0000000000", DifficultyToString(uint64(0xfffffe0000000000)))
	assert.Equal(t, "fffffff800000000", DifficultyToString(uint64(0xfffffff800000000)))
}

func TestIsWorkValid(t *testing.T) {
	// Valid result
	workResult := "205452237a9b01f4"
	hash := "3F93C5CD2E314FA16702189041E68E68C07B27961BF37F0B7705145BEFBA3AA3"

	// Check that we can access these items
	assert.True(t, IsWorkValid(hash, 1, workResult))

	// Invalid result
	workResult = "205452237a9b01f4"
	hash = "3F93C5CD2E314FA16702189041E68E68C07B27961BF37F0B7705145BEFBA3AA3"

	// Check that we can access these items
	assert.False(t, IsWorkValid(hash, 800, workResult))

	// Invalid result
	workResult = "205452237a9b01f4"
	hash = "F1C59E6C738BB82221E082910740BADC58301F8F32291E07CCC4CDBEEAD44348"

	// Check that we can access these items
	assert.False(t, IsWorkValid(hash, 1, workResult))

	hash = "03DDDFF29D3FF3DC41B5374A10A70B49F7AA41E42461511D6A64F346F9C8421E"
	workResult = "00000000002d7708"
	assert.False(t, IsWorkValid(hash, 1, workResult))
}

func TestReverse(t *testing.T) {
	arr := []byte{1, 2, 3, 4, 5}
	reverse(arr)
	assert.Equal(t, []byte{5, 4, 3, 2, 1}, arr)
}

func TestWorkToString(t *testing.T) {
	work := nanopow.NewWork([]byte{1, 2, 3, 4, 5, 6, 7, 8})
	assert.Equal(t, "0102030405060708", WorkToString(work))
}
