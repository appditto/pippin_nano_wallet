package utils

import (
	"errors"
	"os"
	"path"
)

func GetPippinConfigurationRoot() (string, error) {
	homeDir, err := os.UserHomeDir()
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
