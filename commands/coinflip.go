package commands

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/urfave/cli/v3"
	"github.com/v1ctorio/termpet/pet"
)

var CoinflipCommand = &cli.Command{
	Name:   "coinflip",
	Usage:  "Make your pet flip a coin",
	Action: coinfilp,
}

var res [2]string = [2]string{"Head", "Tails"}

func coinfilp(context context.Context, cmd *cli.Command) error {
	pet.Sayln("Flipping the coin right now")
	time.Sleep(time.Second)

	for i := 0; i < 3; i++ {
		fmt.Print(".")
		time.Sleep(time.Second)
	}

	fmt.Println()
	pet.Say("The outcome has been decided, you got %s", res[rand.Intn(2)])
	return nil
}
