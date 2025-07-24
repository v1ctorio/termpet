package commands

import (
	"context"
	"math/rand"

	"github.com/urfave/cli/v3"
	"github.com/v1ctorio/termpet/pet"
)

var normalGreets = [...]string{"Greeting dear user, %s here", "How you doing user? pretty chill here%t", "Why did you decide to name me %s?", "That man page is crazy  ngl%t", "Calling me right now? For real?%t", "I missed you so much!%t"}

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
	pet.Say(normalGreets[rand.Intn(len(normalGreets))], name)

	return nil
}

func RandN(min int, max int) int {
	return rand.Intn(max-min) + min
}
