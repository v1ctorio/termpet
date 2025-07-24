package commands

import (
	"context"
	"fmt"
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
	Name: "init",
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
	},
}

func initPet(ctx context.Context, cmd *cli.Command) (err error) {
	petName := strings.TrimSpace(cmd.StringArg("Pet name"))

	if petName == "" {
		return fmt.Errorf("Please, specify a pet name")
	}
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
			return fmt.Errorf("Pet %s already exists. Choose a different name to init the new pet", petName)
		}

		b, err := tx.CreateBucketIfNotExists(B(petName))
		if err != nil {
			return err
		}
		return b.Put(B("name"), B(petName))
	})
	db.Close()
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

	dbncfg.WriteConfig(dbncfg.ConfigPath, dbncfg.TermpetConfig{
		CommandParser:  dbncfg.Config.CommandParser,
		DatabaseDir:    dbncfg.Config.DatabaseDir,
		DefaultPetName: petName,
	})

	fmt.Printf("Succesfully created pet %s in %s\n", petName, dbncfg.Config.DatabaseDir)
	return nil
}
