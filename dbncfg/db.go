package dbncfg

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/boltdb/bolt"
)

func OpenDB(path string) (*bolt.DB, error) {
	if path == "" {
		path = Config.DatabaseDir
	}

	path, err := SanitizePath(path)
	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
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
