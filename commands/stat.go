package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli/v3"
	"github.com/v1ctorio/termpet/pet"
)

var StatCommand = &cli.Command{
	Name:   "stat",
	Usage:  "View the statistics of your pet",
	Action: stat,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "json",
			Usage: "Prints the output in json to stdout and doesn't pass it through the parser",
			Value: false,
		},
	},
}

func I(str string) int64 {
	i, e := strconv.ParseInt(str, 10, 64)
	if e != nil {
		return 0
	}
	return i
}

func stat(ctx context.Context, cmd *cli.Command) (err error) {
	err = nil

	useJson := cmd.Bool("json")
	p, err := pet.GetPet()
	if err != nil {
		return err
	}

	err = pet.UpdateHunger()
	if err != nil {
		return err
	}

	if !useJson {
		pet.Say("My name is %s", p.Name)
		pet.Say("The latest time you interacted with me was %s, that was %s ", p.LatestInteractionTimestamp, formatUnixTime(I(p.LatestInteractionTimestamp), time.DateTime))
		pet.Say("And my hunger is at %d/48", p.Hunger)
	} else {
		output := fmt.Sprintf("{ name:\"%s\", last_interacted: \"%s\", hunger: \"%d\"}", p.Name, p.LatestInteractionTimestamp, p.Hunger)
		fmt.Println(output)
	}

	return

}
func formatUnixTime(seconds int64, layout string) string {

	location, err := time.LoadLocation("Local")

	if err != nil {
		log.Fatal(err)
	}
	if tz := os.Getenv("TZ"); tz != "" {
		location, err = time.LoadLocation(tz)
		if err != nil {
			log.Fatal(err)
		}
	}
	return time.Unix(seconds, 0).In(location).Format(layout)
}
