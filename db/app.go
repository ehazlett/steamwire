package db

import (
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/ehazlett/steamwire/types"
)

// AddApp adds an application ID to the database
func (d *DB) AddApp(appID string) error {
	// validate
	valid, err := d.IsValidID(appID)
	if err != nil {
		return err
	}

	if !valid {
		return fmt.Errorf("%s is not a valid app id", appID)
	}

	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucketApps))
		return b.Put([]byte(appID), nil)
	})
}

// UpdateAppIndex updates the application ID with the latest
// news GID from an update
func (d *DB) UpdateAppIndex(appID string, index string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucketApps))
		return b.Put([]byte(appID), []byte(index))
	})
}

// DeleteApp deletes an application ID from the database
func (d *DB) DeleteApp(appID string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucketApps))
		return b.Delete([]byte(appID))
	})
}

// GetAppIndex gets the current application index (news GID)
// for the specified application
func (d *DB) GetAppIndex(appID string) (string, error) {
	index := ""
	if err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucketApps))
		index = string(b.Get([]byte(appID)))
		return nil
	}); err != nil {
		return "", err
	}

	return index, nil
}

// GetApps returns all application IDs
func (d *DB) GetApps() ([]string, error) {
	apps := []string{}
	if err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucketApps))
		b.ForEach(func(k, v []byte) error {
			apps = append(apps, string(k))
			return nil
		})
		return nil
	}); err != nil {
		return nil, err
	}

	return apps, nil
}

// GetAppInfo returns the info for the specified application
func (d *DB) GetAppInfo(appID string) (*types.AppInfo, error) {
	info := &types.AppInfo{}
	if err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucketAppList))
		// TODO: improve this lookup; separate table?
		b.ForEach(func(k, v []byte) error {
			if string(v) == appID {
				i, err := getAppInfo(k, v)
				if err != nil {
					return err
				}

				info = i
				return nil
			}

			return nil
		})
		return nil
	}); err != nil {
		return nil, err
	}

	return info, nil
}

func getAppInfo(k, v []byte) (*types.AppInfo, error) {
	appID, err := strconv.Atoi(string(v))
	if err != nil {
		return nil, err
	}
	return &types.AppInfo{
		Name:  string(k),
		AppID: appID,
	}, nil
}
