package dbncfg

import "fmt"

var Config TermpetConfig
var ConfigPath string

func InitConfig() error {
	var err error
	Config, err = readConfig(ConfigPath)
	if err != nil {
		return fmt.Errorf("Error reading or creating the config file in %s. You can change it with the DEFAULT_CONFIG_PATH environment variable. %w", ConfigPath, err)
	}
	return nil
}
