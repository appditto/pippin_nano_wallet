package config

import (
	"os"
	"path"

	"github.com/appditto/pippin_nano_wallet/libs/config/models"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

func ParsePippinConfig() (*models.PippinConfig, error) {
	// Get pippin config path
	pippinConfigPath, err := utils.GetPippinConfigurationRoot()
	if err != nil {
		return nil, err
	}
	// See if file exists
	if _, err := os.Stat(path.Join(pippinConfigPath, "config.yaml")); os.IsNotExist(err) {
		// Create file with defaults
		config := models.PippinConfig{
			Server: models.ServerConfig{},
			Wallet: models.WalletConfig{},
		}
		// Set defaults
		if err := defaults.Set(&config); err != nil {
			return nil, err
		}
		// Write to file
		f, err := os.Create(path.Join(pippinConfigPath, "config.yaml"))
		if err != nil {
			return nil, err
		}
		defer f.Close()
		encoder := yaml.NewEncoder(f)
		defer encoder.Close()
		encoder.SetIndent(2)
		err = encoder.Encode(config)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		// Some other reason we couldn't read the file, maybe permissions
		return nil, err
	}
	// Get fil
	file, err := os.ReadFile(path.Join(pippinConfigPath, "config.yaml"))
	if err != nil {
		return nil, err
	}

	// Parse yaml
	var config models.PippinConfig
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	// Set defaults
	if err := defaults.Set(&config); err != nil {
		return nil, err
	}

	// Validate aspects of the configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}
