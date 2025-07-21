package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/v1ctorio/termpet/commands"
	"github.com/v1ctorio/termpet/dbncfg"
)

// TODO: Add windows and linux support
const DEFAULT_CONFIG_PATH = "~/.config/termpet/termpet.toml"

func main() {

	if dbncfg.ConfigPath = os.Getenv("TERMPET_CONFIG_PATH"); dbncfg.ConfigPath == "" {
		dbncfg.ConfigPath = DEFAULT_CONFIG_PATH
	}

	err := dbncfg.InitConfig()
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
