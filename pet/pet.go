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

type PetData struct {
	Name                       string
	LatestInteractionTimestamp string
	Hunger                     int
}

func (p PetData) save() error {
	err := SetK(PetName, p.Name)
	err = SetK(PetHunger, p.Hunger)
	err = SetK(PetLatestInteractionTimestamp, p.LatestInteractionTimestamp)
	return err
}
func GetPet() (PetData, error) {
	pn, err := GetKNoUpdate(PetName)
	if err != nil {
		return PetData{}, err
	}
	hunger, err := GetKNoUpdate(PetHunger)
	if err != nil {
		return PetData{}, err
	}
	h, err := strconv.Atoi(hunger)
	if err != nil {
		return PetData{}, err
	}
	lit, err := GetKNoUpdate(PetLatestInteractionTimestamp)
	if err != nil {
		return PetData{}, err
	}
	return PetData{
		Name:                       pn,
		LatestInteractionTimestamp: lit,
		Hunger:                     h,
	}, nil

}

const (
	SECOND = 60
	HOUR   = 60
	DAY    = HOUR * 24
	WEEK   = DAY * 7
)

const STARVATION = 24 * 2

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
	return dbncfg.SetV(db, PetLatestInteractionTimestamp.String(), GetCurrentUnixTimestampString())
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
	var value string = ""
	if n, ok := any(val).(int); ok {
		value = strconv.Itoa(n)
	}
	if s, ok := any(val).(string); ok {
		value = s
	}
	db, err := dbncfg.OpenDB(dbncfg.Config.DatabaseDir)
	if err != nil {
		return err
	}
	defer db.Close()

	err = dbncfg.SetV(db, key.String(), value)

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

func GetCurrentUnixTimestampString() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func updateHunger() error {
	pet, err := GetPet()
	if err != nil {
		return err
	}
	hungerToAdd, err := calculateHunger(pet.LatestInteractionTimestamp)
	if err != nil {
		return err
	}
	pet.Hunger = pet.Hunger + hungerToAdd

	err = pet.save()
	return err
}

func calculateHunger(latestInteractionTimestamp string) (int, error) {
	currentTimeStamp := time.Now().Unix()
	LIT, err := strconv.ParseInt(latestInteractionTimestamp, 10, 64)
	if err != nil {
		return 0, err
	}
	difference := currentTimeStamp - LIT

	if difference < HOUR/2 {
		return 0, nil
	}
	if difference < 3*HOUR {
		return 1, nil
	}
	if difference < 3*HOUR {
		return 1, nil
	}
	if difference < 15*HOUR {
		return 15, nil
	}
	if difference < 24*HOUR {
		return 24, nil
	}
	return 35, nil

}
