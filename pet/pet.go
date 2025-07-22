package pet

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/v1ctorio/termpet/dbncfg"
)

var config dbncfg.TermpetConfig = dbncfg.Config

func Say(text string, v ...any) error {

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
