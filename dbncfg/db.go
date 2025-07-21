package dbncfg

import "github.com/boltdb/bolt"

func OpenDB(path string) (*bolt.DB, error) {
	if path == "" {
		path = Config.DatabaseDir
	}

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}
