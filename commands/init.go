package commands

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/urfave/cli/v3"
	"github.com/v1ctorio/termpet/dbncfg"
	"github.com/v1ctorio/termpet/pet"
)

func B(s string) []byte {
	return []byte(s)
}

var InitCommand = &cli.Command{
	Name:      "init",
	UsageText: "Run `termpet init -h` to see the init commands",
	Commands: []*cli.Command{
		{
			Name:  "pet",
			Usage: "Create a new pet",
			Arguments: []cli.Argument{
				&cli.StringArg{
					Name: "Pet name",
				},
			},
			Action: initPet,
		},
		{
			Name:  "startup",
			Usage: "Add the pet to your terminal start script",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "write",
					Usage: "Write the changes to disk instead of just printing them",
					Value: false,
				},
				&cli.StringFlag{
					Name:  "shell",
					Usage: "Specify the desired term for the configuration",
					Value: "bash",
				},
			},
			Action: initStartup,
		},
		{
			Name:  "slack",
			Usage: "Add a slack webhook for the dispatch command",
			Arguments: []cli.Argument{
				&cli.StringArg{
					Name:      "webhook",
					UsageText: "Slack webhook url to use",
				},
			},
			Action: initSlack,
		},
	},
}

func initPet(ctx context.Context, cmd *cli.Command) (err error) {
	petName := strings.TrimSpace(cmd.StringArg("Pet name"))

	if petName == "" {
		return fmt.Errorf("please, specify a pet name")
	}
	dbncfg.PetName = petName
	var db *bolt.DB
	dbncfg.Config.DatabaseDir, err = dbncfg.SanitizePath(dbncfg.Config.DatabaseDir)
	if err != nil {
		return err
	}
	db, err = dbncfg.OpenDB(dbncfg.Config.DatabaseDir)
	if err != nil {
		return err
	}
	println("Creating new pet", petName)
	db.Update(func(tx *bolt.Tx) error {

		bu := tx.Bucket(B(petName))
		if bu != nil {
			return fmt.Errorf("pet %s already exists. Choose a different name to init the new pet", petName)
		}

		b, err := tx.CreateBucketIfNotExists(B(petName))
		if err != nil {
			return err
		}
		return b.Put(B("name"), B(petName))
	})
	db.Close()
	dbncfg.Config, err = dbncfg.WriteConfig(dbncfg.ConfigPath, dbncfg.TermpetConfig{
		CommandParser:  dbncfg.Config.CommandParser,
		DatabaseDir:    dbncfg.Config.DatabaseDir,
		DefaultPetName: petName,
	})
	if err != nil {
		return err
	}
	dbncfg.Config.DefaultPetName = petName

	err = pet.SetK(pet.PetHunger, 0)
	if err != nil {
		return err
	}
	err = pet.SetK(pet.PetLatestInteractionTimestamp, pet.GetCurrentUnixTimestampString())
	if err != nil {
		return err
	}
	err = pet.SetK(pet.PetSickness, "none")
	if err != nil {
		return err
	}

	if dbncfg.Config.DefaultPetName != "" {
		fmt.Printf("Default pet is already set to %s, it will be replaced by %s.\n", dbncfg.Config.DefaultPetName, petName)
	}

	fmt.Printf("Succesfully created pet %s in %s\n", petName, dbncfg.Config.DatabaseDir)
	return nil
}

func initStartup(ctx context.Context, cmd *cli.Command) (err error) {
	var write = cmd.Bool("write")

	if write {

		if e, err := exec.LookPath("termpet"); e == "" || err != nil {
			return fmt.Errorf("termpet not found in path, add it first to make startup work %w", err)
		}
	}

	var shell string
	shell = cmd.String("shell")

	if runtime.GOOS == "windows" {
		shell = "powershell"
	} else if strings.Contains(os.Getenv("SHELL"), "fish") {
		shell = "fish"
	} else if strings.Contains(os.Getenv("SHELL"), "zsh") {
		shell = "zsh"
	}

	if shell == "bash" {
		println("Shell assumed to be bash. If using zsh or fish, specify so with the flag `--shell`")
	}

	stringToAdd := `termpet greet`

	_ = stringToAdd

	var fileToEdit string
	switch shell {
	case "powershell":
		out, err := exec.Command("powershell", "-NoProfile", "-Command", "Write-Output $PROFILE").CombinedOutput()
		if err != nil {
			return err
		}
		fileToEdit = strings.TrimSpace(string(out))
	case "bash":
		fileToEdit, err = dbncfg.SanitizePath("~/.bashrc")
		if err != nil {
			return err
		}
	case "zsh":
		fileToEdit, err = dbncfg.SanitizePath("~/.bashrc")
		if err != nil {
			return err
		}
	case "fish":

		fileToEdit, err = dbncfg.SanitizePath("&/fish/config.fish")
		if err != nil {
			return err
		}
	}

	if write {

		f, err := os.OpenFile(fileToEdit, os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		_, err = f.WriteString("\n" + stringToAdd)

		if err != nil {
			return fmt.Errorf("Error writing to the shell config file %w", err)
		}

		pet.YellowLn("Succefully added pet to shell startup")
	} else {
		fmt.Printf("To make your terminal pet appear on startup add `%s` to %s\n", stringToAdd, fileToEdit)
	}

	return nil
}

func initSlack(ctx context.Context, cmd *cli.Command) (err error) {
	webhook := strings.TrimSpace(cmd.StringArg("webhook"))
	if webhook == "" {
		return fmt.Errorf("Provide a slack webhook as an argument!")
	}
	if _, err := url.ParseRequestURI(webhook); err != nil {

		return fmt.Errorf("Invalid url provided %w", err)
	}

	err = pet.SetK(pet.SlackWebhook, webhook)
	pet.YellowLn("Webhook url succesfully saved")
	return err
}
