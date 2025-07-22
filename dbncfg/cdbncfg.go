package dbncfg

import (
	"fmt"
	"os"
)

var Config TermpetConfig
var ConfigPath string
var PetName string

func InitConfig() error {
	var err error
	Config, err = readConfig(ConfigPath)
	if err != nil {
		return fmt.Errorf("Error reading or creating the config file in %s. You can change it with the DEFAULT_CONFIG_PATH environment variable. %w", ConfigPath, err)
	}

	if pN := os.Getenv("TERMPET_PET"); pN != "" {
		PetName = pN
	} else {
		PetName = Config.DefaultPetName
	}

	return nil
}
