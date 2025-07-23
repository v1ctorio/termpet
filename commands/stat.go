package commands

import (
	"context"
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
	name, err := pet.GetName()
	if err != nil {
		return err
	}
	lit, err := pet.GetKNoUpdate(pet.PetLatestInteractionTimestamp)

	if !useJson {
		pet.Say("My name is %s", name)
		pet.Say("The latest time you interacted with me was %s, that was %s ", lit, formatUnixTime(I(lit), time.DateTime))
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
