package commands

import (
	"context"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/urfave/cli/v3"
	"github.com/v1ctorio/termpet/dbncfg"
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
	var db *bolt.DB
	db, err = dbncfg.OpenDB(dbncfg.Config.DatabaseDir)
	if err != nil {
		return err
	}
	defer db.Close()
	println("Creating new pet", petName)
	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(B(petName))
		if err != nil {
			return err
		}
		return b.Put(B("name"), B(petName))
	})
	return nil
}
