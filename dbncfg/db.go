package dbncfg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boltdb/bolt"
)

var noBucketFound = "No bucket %s found. Please init your pet with the `init pet` subcommand"

func B(s string) []byte {
	return []byte(s)
}

func OpenDB(path string) (*bolt.DB, error) {
	if path == "" {
		path = Config.DatabaseDir
	}

	path, err := SanitizePath(path)
	path, err = SanitizePath(path)
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

func GetV(db *bolt.DB, key string) (string, error) {

	if PetName == "" {
		return "", fmt.Errorf("Error trying to read database. No pet name provided")
	}

	var value string

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(B(PetName))
		if bucket == nil {
			return fmt.Errorf(noBucketFound, PetName)
		}
		v := bucket.Get(B(key))
		if v != nil {
			value = string(v)
		}
		return nil
	})

	if err != nil {
		return "", err
	}
	return value, nil
}

func SetV(db *bolt.DB, key string, value string) error {
	if key == "" {
		return fmt.Errorf("No key provided for errorf")
	}

	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(B(PetName))
		if bucket == nil {
			return fmt.Errorf(noBucketFound, PetName)
		}
		return bucket.Put(B(key), B(value))
	})
	db.Close()
	return err
}
