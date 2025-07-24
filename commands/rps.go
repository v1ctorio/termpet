package commands

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v3"
	"github.com/v1ctorio/termpet/pet"
)

var RpsCommand = &cli.Command{
	Name:   "rps",
	Usage:  "Rock, papaer, scissors",
	Action: rps,
}

var hands = [3]string{"Rock", "Paper", "Scissors"}

func rps(ctx context.Context, cmd *cli.Command) error {

	name, err := pet.GetName()
	if err != nil {
		return err
	}
	pet.Sayln("Let's play rock paper scissors")

	for i := 0; i < 3; i++ {
		fmt.Print(".")
		time.Sleep(time.Second)
	}
	fmt.Println("\nGo!")

	petChosen := hands[rand.Intn(3)]
	time.Sleep(time.Second / 2)
	pet.YellowLn("0: Rock; 1: Paper; 2: Scissors")
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	i, err := strconv.Atoi(input)
	if err != nil {
		return err
	}

	input = hands[i]

	// For transparency. The switch below is AI generated. There was no point on making such thing thats been made a thousand times again

	fmt.Printf("%s chose %s\n", name, petChosen)

	switch {
	case input == petChosen:
		pet.Sayln("It's a tie! 🤝")
	case (input == "Rock" && petChosen == "Scissors") ||
		(input == "Paper" && petChosen == "Rock") ||
		(input == "Scissors" && petChosen == "Paper"):
		pet.Sayln("You win! Good job!")
	default:
		pet.Sayln("I win this time! Better luck next round! 😸")
	}
	if err != nil {
		return err
	}

	return nil
}
