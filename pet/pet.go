package pet

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/v1ctorio/termpet/dbncfg"
)

type petValidKey string

const (
	PetName                       petValidKey = "name"
	PetLatestInteractionTimestamp petValidKey = "latestinteractiontime"
	PetHunger                     petValidKey = "hunger"
	PetSickness                   petValidKey = "sickness"
)

// I know this is hardcoded and kinda trashy but idk how to make it better without dependencies
func (k petValidKey) String() string {
	return string(k)
}

type PetData struct {
	Name                       string
	LatestInteractionTimestamp string
	Hunger                     int
	Sickness                   string
}

func (p PetData) Save() error {

	db, err := dbncfg.OpenDB(config.DatabaseDir)
	if err != nil {
		return err
	}
	err = dbncfg.SetV(db, PetName.String(), p.Name)
	err = dbncfg.SetV(db, PetHunger.String(), p.Hunger)
	err = dbncfg.SetV(db, PetLatestInteractionTimestamp.String(), p.LatestInteractionTimestamp)
	err = dbncfg.SetV(db, PetSickness.String(), p.Sickness)
	db.Close()
	return err
}
func (p *PetData) UpdateLatestInteractionTime() {
	p.LatestInteractionTimestamp = strconv.FormatInt(time.Now().Unix(), 10)
}
func GetPet() (pd PetData, err error) {
	pd = PetData{}
	pn, err := GetKNoUpdate(PetName)
	if err != nil {
		return
	}
	hunger, err := GetKNoUpdate(PetHunger)
	if err != nil {
		return
	}
	h, err := strconv.Atoi(hunger)
	if err != nil {
		return
	}
	lit, err := GetKNoUpdate(PetLatestInteractionTimestamp)
	if err != nil {
		return
	}
	s, err := GetKNoUpdate(PetSickness)
	if err != nil {
		return
	}
	return PetData{
		Name:                       pn,
		LatestInteractionTimestamp: lit,
		Hunger:                     h,
		Sickness:                   s,
	}, nil

}

const (
	SECOND = 1
	MINUTE = 60 * SECOND
	HOUR   = MINUTE * 60
	DAY    = HOUR * 24
	WEEK   = DAY * 7
)

const STARVATION = 24 * 2

var config *dbncfg.TermpetConfig = &dbncfg.Config

var SayContent string = ""

func Sayln(text string, v ...any) error {
	if strings.Contains(text, "%t") {
		v = []any{}
		text = strings.ReplaceAll(text, "%t", "")
	}
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

	if strings.Contains(text, "%t") {
		v = []any{}
		text = strings.ReplaceAll(text, "%t", "")
	}

	if SayContent == "" {
		SayContent = fmt.Sprintf(fmt.Sprintf(text, v...))
	} else {
		SayContent = SayContent + "\n" + fmt.Sprintf(text, v...)
	}
}
func UpdateLatestInteractionTime() error {
	db, err := dbncfg.OpenDB(dbncfg.Config.DatabaseDir)
	if err != nil {
		return err
	}
	err = dbncfg.SetV(db, PetLatestInteractionTimestamp.String(), GetCurrentUnixTimestampString())
	db.Close()
	return err

}

func GetName() (name string, err error) {
	err = nil
	db, err := dbncfg.OpenDB(dbncfg.Config.DatabaseDir)
	if err != nil {
		return
	}
	name, err = dbncfg.GetV(db, PetName.String())

	db.Close()

	err = UpdateLatestInteractionTime()
	if err != nil {
		return
	}

	return
}

func YellowLn(text string, v ...any) {
	fmt.Print("\033[93;1;3m" + fmt.Sprintf(text, v...) + "\033[0m\n")
}

func SetK[T string | int](key petValidKey, val T) error {

	db, err := dbncfg.OpenDB(dbncfg.Config.DatabaseDir)
	if err != nil {
		return err
	}

	err = dbncfg.SetV(db, key.String(), val)
	db.Close()

	return err

}

func getKey(key petValidKey, doUpdateLatestInteractionTime bool) (value string, err error) {
	err = nil
	db, err := dbncfg.OpenDB(dbncfg.Config.DatabaseDir)
	if err != nil {
		return
	}

	value, err = dbncfg.GetV(db, key.String())
	db.Close()
	if doUpdateLatestInteractionTime {
		err = UpdateLatestInteractionTime()
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

func UpdateHunger() error {
	pet, err := GetPet()
	if err != nil {
		return err
	}
	hungerToAdd, err := calculateHunger(pet.LatestInteractionTimestamp)
	fmt.Printf("Adding %d hunger to the pet\n", hungerToAdd)
	if err != nil {
		return err
	}
	pet.Hunger = pet.Hunger + hungerToAdd

	if pet.Hunger > 18 {
		YellowLn("%s is starving! try to feed them.", pet.Name)
	}

	if pet.Hunger > 24 && pet.Sickness != "hungry" {
		YellowLn("%s is so hungry that they got sick!", pet.Name)
		pet.Sickness = "hungry"
	}

	err = pet.Save()
	return err
}

func calculateHunger(latestInteractionTimestamp string) (int, error) {
	currentTimeStamp := time.Now().Unix()
	LIT, err := strconv.ParseInt(latestInteractionTimestamp, 10, 64)
	if err != nil {
		return 0, err
	}
	difference := currentTimeStamp - LIT
	//println(difference)

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
