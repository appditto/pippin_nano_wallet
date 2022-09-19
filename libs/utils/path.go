package utils

import (
	"errors"
	"os"
	"path"
)

func GetPippinConfigurationRoot() (string, error) {
	var homeDir string
	var err error
	if GetEnv("PIPPIN_HOME", "") == "" {
		homeDir, err = os.UserHomeDir()
	} else {
		homeDir = GetEnv("PIPPIN_HOME", "")
	}
	if err != nil {
		return "", err
	}

	pippinDataDir := path.Join(homeDir, "PippinData")
	if _, err := os.Stat(pippinDataDir); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(pippinDataDir, os.ModePerm); err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	return pippinDataDir, nil
}
