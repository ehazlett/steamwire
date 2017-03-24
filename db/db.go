package db

import (
	"fmt"

	"github.com/boltdb/bolt"
)

const (
	dbBucketApps    = "apps"
	dbBucketAppList = "applist"
	dbBucketMeta    = "meta"
)

// DB is the steamwire datastore
type DB struct {
	db *bolt.DB
}

// NewDB returns a new steamwire datastore
func NewDB(db *bolt.DB) (*DB, error) {
	// ensure bucket is created
	for _, b := range []string{dbBucketApps, dbBucketAppList, dbBucketMeta} {
		if err := db.Update(func(tx *bolt.Tx) error {
			if _, err := tx.CreateBucketIfNotExists([]byte(b)); err != nil {
				return fmt.Errorf("error creating bucket: %s", err)
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}

	return &DB{
		db: db,
	}, nil
}
