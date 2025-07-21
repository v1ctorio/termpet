package dbncfg

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
)

type TermpetConfig struct {
	DatabaseDir   string //Where the pet db file is stored
	CommandParser string //Cmmand
}

const DEFAULT_CONFIG_PATH = "~/.config/termpet/termpet.toml"
const DEFAULT_PET_DIR = "~/.config/termpet/pet.db"
const DEFAULT_COMMAND_PARSER = "cowsay -f koala \"{}\""

func readConfig(dir string) (cfg TermpetConfig, err error) {

	cfg = TermpetConfig{}
	err = nil

	if dir == "" {
		dir = DEFAULT_CONFIG_PATH
	}
	txt, err := os.ReadFile(dir)
	if err != nil {
		println("No config found. Initializing default config in ", dir)
		cfg, err = initConfig(dir, DEFAULT_PET_DIR, DEFAULT_CONFIG_PATH)
		if err != nil {
			return
		}
	} else {
		err = toml.Unmarshal(txt, &cfg)
		if err != nil {
			return
		}
	}

	return

}

func initConfig(path string, dbDir string, commandParser string) (TermpetConfig, error) {

	cfg := TermpetConfig{
		CommandParser: commandParser,
		DatabaseDir:   dbDir,
	}
	tml, err := toml.Marshal(cfg)
	if err != nil {
		return TermpetConfig{}, fmt.Errorf("Error %encoding the config in toml %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()
	if err != nil {
		return TermpetConfig{}, err
	}

	_, err = f.Write(tml)
	if err != nil {
		return TermpetConfig{}, err
	}
	f.Sync()
	return cfg, nil
}
