package pet

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/v1ctorio/termpet/dbncfg"
)

const (
	PetName                       string = "name"
	PetLatestInteractionTimestamp        = "latestinteractiontimestamp"
)

var config *dbncfg.TermpetConfig = &dbncfg.Config

var SayContent string = ""

func Sayln(text string, v ...any) error {

	formatted := fmt.Sprintf(text, v...)

	cmdParts := strings.Fields(strings.Replace(config.CommandParser, "{}", formatted, 1))
	if len(cmdParts) == 0 {
		return fmt.Errorf("No command to execute")
	}
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)

	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error executing the parser command %w", err)
	}
	return nil
}

// Adds each say to a bucket of say things to print with sayln on exit
func Say(text string, v ...any) {
	SayContent = SayContent + "\n" + fmt.Sprintf(text, v...)
}
func updateLatestInteractionTime(db *bolt.DB) error {
	return dbncfg.SetV(db, PetLatestInteractionTimestamp, getCurrentUnixTimestampString())
}

func GetName() (name string, err error) {
	err = nil
	db, err := dbncfg.OpenDB(dbncfg.Config.DatabaseDir)
	if err != nil {
		return
	}
	defer db.Close()

	err = updateLatestInteractionTime(db)
	if err != nil {
		return
	}

	name, err = dbncfg.GetV(db, PetName)

	return

}

func getCurrentUnixTimestampString() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
