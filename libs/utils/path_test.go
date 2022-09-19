package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPippinConfigurationRoot(t *testing.T) {
	defer os.RemoveAll(".testdata")
	// Set env
	os.Setenv("HOME", ".testdata")
	defer os.Unsetenv("HOME")

	// Get home dir
	pippinDataDir, err := GetPippinConfigurationRoot()
	assert.Equal(t, err, nil)

	// Check if dir exists
	_, err = os.Stat(pippinDataDir)
	assert.Equal(t, nil, nil)

	assert.Equal(t, ".testdata/PippinData", pippinDataDir)
}

func TestGetPippinConfigurationRootPippinHome(t *testing.T) {
	defer os.RemoveAll(".testpippinhome")
	// Set env
	os.Setenv("PIPPIN_HOME", ".testpippinhome")
	defer os.Unsetenv("PIPPIN_HOME")

	// Get home dir
	pippinDataDir, err := GetPippinConfigurationRoot()
	assert.Equal(t, err, nil)

	// Check if dir exists
	_, err = os.Stat(pippinDataDir)
	assert.Equal(t, nil, nil)

	assert.Equal(t, ".testpippinhome/PippinData", pippinDataDir)
}
