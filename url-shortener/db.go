package urlshort

import (
	"fmt"

	"github.com/boltdb/bolt"
)

const BUCKET_NAME = "paths"

func CreatePath(db *bolt.DB, path string, url string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(BUCKET_NAME))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		err = b.Put([]byte(path), []byte(url))
		return err
	})
}

func GetPath(db *bolt.DB, path string, value *string) error {
	return db.View(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(BUCKET_NAME))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		v := b.Get([]byte(path))
		*value = string(v)
		return nil
	})
}
