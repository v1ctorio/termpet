package commands

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	_ "embed"

	"github.com/urfave/cli/v3"
	"github.com/v1ctorio/termpet/pet"
)

var DispatchCommand = &cli.Command{
	Name:   "dispatch",
	Usage:  "Dispatch a message to a service. Only Slack supported atm",
	Action: dispatch,
}

//go:embed assets/slacktemplate.json
var dispatchSlackContent string

func dispatch(ctx context.Context, cmd *cli.Command) error {

	name, err := pet.GetName()
	if err != nil {
		return err
	}
	webhook, err := pet.GetK(pet.SlackWebhook)
	if err != nil {
		return err
	}

	var message string
	for i := 0; i < cmd.Args().Len(); i++ {
		if message == "" {
			message = cmd.Args().Get(i)
		} else {
			message = message + " " + cmd.Args().Get(i)

		}
	}

	parsed, err := pet.ParseWithC(message)
	formattedPayload := fmt.Sprintf(dispatchSlackContent, name, parsed)

	req, err := http.NewRequest("POST", webhook, strings.NewReader(formattedPayload))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Error dispatching the message to slack. Check if your webhook is valid.")
	}

	return nil

}
