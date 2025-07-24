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

const PERMS = 0666
const DEFAULT_CONFIG_PATH = "&/termpet/termpet.toml"
const DEFAULT_PET_DB_PATH = "&/termpet/pet.db"
const DEFAULT_COMMAND_PARSER = "cowsay -f koala \"{}\""

func readConfig(dir string) (cfg TermpetConfig, err error) {

	cfg = TermpetConfig{}
	err = nil

	if dir == "" {
		dir, err = SanitizePath(DEFAULT_CONFIG_PATH)
		if err != nil {
			return
		}
	}
	//println("Reading config from ", dir)
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

	tml, err := toml.Marshal(cfg)
	if err != nil {
		return TermpetConfig{}, fmt.Errorf("Error encoding the config in toml %w", err)
	}

	err = CreateDirForFile(path)
	if err != nil {
		return TermpetConfig{}, err
	}
	if _, err := os.Stat(path); err == nil {

		err = os.Remove(path)
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

func CreateDirForFile(fpath string) error {
	dir := filepath.Dir(fpath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}
