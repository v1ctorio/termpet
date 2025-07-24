package commands

import (
	"context"

	"github.com/urfave/cli/v3"
)

var FeedCommand = &cli.Command{
	Name:   "feed",
	Usage:  "Feed the pet",
	Action: feed,
}

func feed(ctx context.Context, cmd *cli.Command) error {
	return nil
}
