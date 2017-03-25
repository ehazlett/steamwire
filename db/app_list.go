package db

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/ehazlett/steamwire/types"
	"github.com/sirupsen/logrus"
)

// UpdateAppList updates the local cache of steam apps
func (d *DB) UpdateAppList(list []*types.AppInfo) error {
	if err := d.db.Update(func(tx *bolt.Tx) error {
		for _, app := range list {
			b := tx.Bucket([]byte(dbBucketAppList))
			id := strconv.Itoa(app.AppID)
			if err := b.Put([]byte(app.Name), []byte(id)); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	// update last updated meta
	if err := d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucketMeta))
		timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
		return b.Put([]byte("updated"), []byte(timestamp))
	}); err != nil {
		return err
	}
	return nil
}

// GetAppListLastUpdated returns the `time.Time` of the last update
func (d *DB) GetAppListLastUpdated() (time.Time, error) {
	var updated time.Time
	if err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucketMeta))
		s := string(b.Get([]byte("updated")))
		// check for empty
		if s == "" {
			return nil
		}
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		updated = time.Unix(i, 0)
		return nil
	}); err != nil {
		return updated, err
	}

	return updated, nil
}

// FindApp returns any app that matches a simple prefix
func (d *DB) FindApp(prefix string) ([]*types.AppInfo, error) {
	apps := []*types.AppInfo{}

	if err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucketAppList))

		logrus.WithFields(logrus.Fields{
			"prefix": prefix,
		}).Debug("searching app list")

		b.ForEach(func(k, v []byte) error {
			if strings.Contains(strings.ToLower(string(k)), strings.ToLower(prefix)) {
				info, err := getAppInfo(k, v)
				if err != nil {
					return err
				}
				apps = append(apps, info)
			}

			return nil
		})
		logrus.WithFields(logrus.Fields{
			"numOfApps": len(apps),
		}).Debug("search results")

		return nil
	}); err != nil {
		return nil, err
	}
	return apps, nil
}

// GetAppList returns the local app list cache
func (d *DB) GetAppList() ([]*types.AppInfo, error) {
	apps := []*types.AppInfo{}
	if err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucketApps))
		b.ForEach(func(k, v []byte) error {
			info, err := getAppInfo(k, v)
			if err != nil {
				return err
			}
			apps = append(apps, info)
			return nil
		})
		return nil
	}); err != nil {
		return nil, err
	}

	return apps, nil
}

// IsValidID returns whether the specified application id is known or not
func (d *DB) IsValidID(appID string) (bool, error) {
	valid := false
	if err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucketAppList))
		// TODO: optimize; separate bucket?
		b.ForEach(func(k, v []byte) error {
			if string(v) == appID {
				valid = true
				return nil
			}
			return nil
		})
		return nil
	}); err != nil {
		return valid, err
	}

	return valid, nil
}
