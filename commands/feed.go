package commands

import (
	"context"
	"strconv"

	"github.com/urfave/cli/v3"
	"github.com/v1ctorio/termpet/pet"
)

var FeedCommand = &cli.Command{
	Name:   "feed",
	Usage:  "Feed the pet",
	Action: feed,
}

func feed(ctx context.Context, cmd *cli.Command) error {
	u, err := pet.GetPet()
	if err != nil {
		return err
	}
	//println(u.LatestInteractionTimestamp)
	lit, err := strconv.Atoi(u.LatestInteractionTimestamp)

	if err != nil {
		return err
	}

	if lit < pet.MINUTE {
		pet.YellowLn("You need to wait at least one minute before feeding %s again!", u.Name)
		return nil
	}

	err = pet.UpdateHunger()

	err = pet.UpdateLatestInteractionTime()
	if err != nil {
		return err
	}

	u.Hunger = u.Hunger - 1

	if u.Hunger <= 0 {
		u.Hunger = 0
		pet.Say("I'm not hungry at all, user!")
	}
	if u.Hunger > 48 {
		u.Hunger = 46
	}

	if u.Sickness == "hunger" && u.Hunger <= 5 {
		u.Sickness = "none"
		pet.YellowLn("%s recovered from their sickness!", u.Name)
	}

	pet.YellowLn("%s hunger is now at %d/48", u.Name, u.Hunger)
	err = u.Save()
	return err
}
