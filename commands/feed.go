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
	lit, err := strconv.Atoi(u.LatestInteractionTimestamp)
	_ = lit
	if err != nil {
		return err
	}
	err = pet.UpdateHunger()

	err = pet.UpdateLatestInteractionTime()
	if err != nil {
		return err
	}
	return nil
}
