package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println(SanitizePath("&/termpet/pet.db"))
	fmt.Println(SanitizePath("&/termpet/termpet.toml"))
}

func SanitizePath(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	configDir, err := os.UserConfigDir()

	if err != nil {
		return "", err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path = strings.Replace(path, "&", configDir, 1)
	path = strings.Replace(path, "~", homeDir, 1)
	path = filepath.Clean(path)

	return path, nil
}
