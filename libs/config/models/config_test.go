package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRandomRep(t *testing.T) {
	config := PippinConfig{}

	_, err := config.GetRandomRep()
	assert.ErrorIs(t, err, ErrNoRepsConfigured)

	config.Wallet.Banano = false
	config.Wallet.PreconfiguredRepresentativesNano = []string{"nano_1"}
	rep, err := config.GetRandomRep()
	assert.Nil(t, err)
	assert.Equal(t, "nano_1", rep)

	config.Wallet.Banano = true
	config.Wallet.PreconfiguredRepresentativesBanano = []string{"ban_1"}
	rep, err = config.GetRandomRep()
	assert.Nil(t, err)
	assert.Equal(t, "ban_1", rep)
}
