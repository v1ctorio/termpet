package dbncfg

import "fmt"

var Config TermpetConfig
var ConfigPath string

func init() {
	var err error
	Config, err = readConfig(ConfigPath)
	if err != nil {
		fmt.Errorf("Error reading or creating the config file in %s. You can change it with the DEFAULT_CONFIG_PATH enviroment var. %w", ConfigPath, err)
	}

}
