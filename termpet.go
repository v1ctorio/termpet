package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/urfave/cli/v3"
	"github.com/v1ctorio/termpet/commands"
	"github.com/v1ctorio/termpet/dbncfg"
	"github.com/v1ctorio/termpet/pet"
)

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
			commands.FeedCommand,
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Print stack traces",
				Value: false,
			},
		},
		Action: noCommand,
	}

	err = cmd.Run(context.Background(), os.Args)
	if err != nil {
		doDebug := cmd.Bool("debug")
		if doDebug {
			fmt.Fprintf(os.Stderr, "%+v.\n", err)
			debug.PrintStack()

		} else {
			fmt.Fprintf(os.Stderr, "%+v.\n", err)
		}

		os.Exit(1)
	}

	if pet.SayContent != "" {
		err = pet.Sayln("%s", pet.SayContent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v.\n", err)
		}
		os.Exit(0)
	}

}

func noCommand(ctx context.Context, cmd *cli.Command) error {
	var keepListening bool = true
	verbs := []string{"p"}

	reader := bufio.NewReader(os.Stdin)
	name, err := pet.GetName()
	if err != nil {
		return err
	}
	pet.Sayln(commands.NormalGreets[rand.Intn(len(commands.NormalGreets))], name)
	fmt.Println("p: pet; f: feed")

	for keepListening {

		input, err := reader.ReadString('\n')

		input = strings.TrimSpace(strings.ReplaceAll(input, "\n", ""))
		if err != nil {
			return err
		}

		if input == "p" {
			err := pet.UpdateLatestInteractionTime()
			err = pet.UpdateHunger()
			if err != nil {
				return err
			}
			pet.YellowLn("You pet %s\n", name)
		}
		if input == "f" {
			keepListening = false
			err = commands.FeedCommand.Run(context.Background(), []string{})
			if err != nil {
				return err
			}
		}

		if !slices.Contains(verbs, string(input)) {
			keepListening = false
		}
		_ = input
	}
	return nil
}
