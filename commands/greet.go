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
	name, err := pet.GetName()
	if err != nil {
		return
	}
	err = pet.UpdateHunger()
	if err != nil {
		return err
	}
	pet.Say("Greetings dear user, %s here", name)

	return nil
}
