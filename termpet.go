package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
	"github.com/v1ctorio/termpet/commands"
	"github.com/v1ctorio/termpet/dbncfg"
)

// TODO: Add windows and linux support
const DEFAULT_CONFIG_PATH = "~/.config/termpet/termpet.toml"

func main() {
	var err error

	if dbncfg.ConfigPath = os.Getenv("TERMPET_CONFIG_PATH"); dbncfg.ConfigPath == "" {
		dbncfg.ConfigPath = DEFAULT_CONFIG_PATH
	}
	dbncfg.ConfigPath, err = sanitizePath(dbncfg.ConfigPath)

	err = dbncfg.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	cmd := &cli.Command{
		Name:  "Termpet",
		Usage: "Take care of your pet!",
		Commands: []*cli.Command{
			commands.InitCommand,
			commands.GreetCommand,
		},
	}

	err = cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func sanitizePath(path string) (string, error) {

	configDir, err := os.UserConfigDir()

	if err != nil {
		return "", err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path = strings.Replace(path, "&", configDir, 1)
	path = strings.Replace(path, "~", homeDir, 1)
	path = filepath.Clean(path)

	return path, nil
}
