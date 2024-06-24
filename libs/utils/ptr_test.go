package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToPtr(t *testing.T) {
	asPtr := ToPtr("test")
	assert.Equal(t, "test", *asPtr)
	asPtrUi := ToPtr[uint](1)
	assert.Equal(t, uint(1), *asPtrUi)
}
