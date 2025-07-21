package dbncfg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type TermpetConfig struct {
	DatabaseDir    string //Where the pet db file is stored
	CommandParser  string //Cmmand
	DefaultPetName string
}

const PERMS = 0644
const DEFAULT_CONFIG_PATH = "~/.config/termpet/termpet.toml"
const DEFAULT_PET_DB_PATH = "~/.config/termpet/pet.db"
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
		cfg, err = WriteConfig(dir, TermpetConfig{
			DatabaseDir:   DEFAULT_PET_DB_PATH,
			CommandParser: DEFAULT_COMMAND_PARSER,
		})
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

func WriteConfig(path string, cfg TermpetConfig) (TermpetConfig, error) {

	path, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return TermpetConfig{}, err
	}

	tml, err := toml.Marshal(cfg)
	if err != nil {
		return TermpetConfig{}, fmt.Errorf("Error encoding the config in toml %w", err)
	}

	d := filepath.Dir(path)
	if _, err := os.Stat(d); os.IsNotExist(err) {
		err := os.MkdirAll(d, PERMS)
		if err != nil {
			return TermpetConfig{}, err
		}
	} else {
		if err != nil {
			return TermpetConfig{}, err
		}
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, PERMS)
	if err != nil {
		return TermpetConfig{}, err
	}
	defer f.Close()

	_, err = f.Write(tml)
	if err != nil {
		return TermpetConfig{}, err
	}
	f.Sync()
	return cfg, nil
}
