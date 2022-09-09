package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate64HexHash(t *testing.T) {
	valid := "1A2E95A2DCF03143297572EAEC496F6913D5001D2F28A728B35CB274294D5A14"
	assert.Equal(t, true, Validate64HexHash(valid))
	// invalid, not hex
	invalid := "1A2E95A2DCZ03143297572EAEC496F6913D5001D2F28A728B35CB274294D5A14"
	assert.Equal(t, false, Validate64HexHash(invalid))
	invalid = "1A2E95A2DCF03143297572EAEC496F6913D5001D2F28A728B35CB274294D5A1"
	assert.Equal(t, false, Validate64HexHash(invalid))
}
