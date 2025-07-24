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

type petValidKey string

const (
	PetName                       petValidKey = "name"
	PetLatestInteractionTimestamp petValidKey = "latestinteractiontime"
	PetHunger                     petValidKey = "hunger"
)

// I know this is hardcoded and kinda trashy but idk how to make it better without dependencies
func (k petValidKey) String() string {
	return string(k)
}

const (
	SECOND = 60
	DAY    = SECOND * 24
	WEEK   = DAY * 7
)

var config *dbncfg.TermpetConfig = &dbncfg.Config

var SayContent string = ""

func Sayln(text string, v ...any) error {

	formatted := fmt.Sprintf(text, v...)

	cmdParts := strings.Fields(strings.Replace(config.CommandParser, "{}", formatted, 1))
	if len(cmdParts) == 0 {
		return fmt.Errorf("No command to execute")
	}

	executable, err := exec.LookPath(cmdParts[0])
	if err != nil || executable == "" {
		return fmt.Errorf("No executable %s found. Check if it's in path", cmdParts[0])
	}
	cmd := exec.Command(executable, cmdParts[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
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
	return dbncfg.SetV(db, PetLatestInteractionTimestamp.String(), getCurrentUnixTimestampString())
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

	name, err = dbncfg.GetV(db, PetName.String())

	return
}

func YellowLn(text string, v ...any) {
	fmt.Print("\033[93;1;3m" + fmt.Sprintf(text, v...) + "\033[0m")
}

func SetK[T string | int](key petValidKey, val T) error {
	db, err := dbncfg.OpenDB(dbncfg.Config.DatabaseDir)
	if err != nil {
		return err
	}
	defer db.Close()

	err = dbncfg.SetV(db, key.String(), val)

	return err

}

func getKey(key petValidKey, doUpdateLatestInteractionTime bool) (value string, err error) {
	err = nil
	db, err := dbncfg.OpenDB(dbncfg.Config.DatabaseDir)
	if err != nil {
		return
	}
	defer db.Close()

	value, err = dbncfg.GetV(db, key.String())

	if doUpdateLatestInteractionTime {
		err = updateLatestInteractionTime(db)
	}

	if err != nil {
		return "", err
	}
	return
}

func GetK(key petValidKey) (string, error) {
	return getKey(key, true)
}

func GetKNoUpdate(key petValidKey) (string, error) {
	return getKey(key, false)
}

func getCurrentUnixTimestampString() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
