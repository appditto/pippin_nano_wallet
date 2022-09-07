package config

import (
	"github.com/appditto/pippin_nano_wallet/libs/config/models"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
)

func ParsePippinConfig() (*models.PippinConfig, error) {
	// Get pippin config path
	pippinConfigPath, err := utils.GetPippinConfigurationRoot()
	if err != nil {
		return nil, err
	}

}
