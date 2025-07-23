package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/v1ctorio/termpet/commands"
	"github.com/v1ctorio/termpet/dbncfg"
	"github.com/v1ctorio/termpet/pet"
)

// TODO: Add windows and linux support

func main() {
	var err error

	if dbncfg.ConfigPath = os.Getenv("TERMPET_CONFIG_PATH"); dbncfg.ConfigPath == "" {

		dbncfg.ConfigPath, err = dbncfg.SanitizePath(dbncfg.DEFAULT_CONFIG_PATH)
		if err != nil {
			log.Fatal(err)
		}
	}
	dbncfg.ConfigPath, err = dbncfg.SanitizePath(dbncfg.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

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
			commands.StatCommand,
		},
	}

	err = cmd.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v.\n", err)
		os.Exit(1)
	}

	if pet.SayContent != "" {
		err = pet.Sayln("%s", pet.SayContent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v.\n", err)
		}
		os.Exit(0)
	}

}
