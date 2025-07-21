package commands

import (
	"context"

	"github.com/urfave/cli/v3"
	"github.com/v1ctorio/termpet/pet"
)

var GreetCommand = &cli.Command{
	Name:   "greet",
	Usage:  "Greet your terminal pet!",
	Action: greet,
}

func greet(ctx context.Context, cmd *cli.Command) (err error) {
	println("Greet command called")
	err = pet.Say("Greetings dear user")
	if err != nil {
		return err
	}
	return nil
}
