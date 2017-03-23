package server

import "github.com/boltdb/bolt"

func (s *Server) AddApp(appID string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucketName))
		return b.Put([]byte(appID), nil)
	})
}

func (s *Server) UpdateAppIndex(appID string, index string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucketName))
		return b.Put([]byte(appID), []byte(index))
	})
}
func (s *Server) DeleteApp(appID string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucketName))
		return b.Delete([]byte(appID))
	})
}

func (s *Server) GetAppIndex(appID string) (string, error) {
	index := ""
	if err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucketName))
		index = string(b.Get([]byte(appID)))
		return nil
	}); err != nil {
		return "", err
	}

	return index, nil
}

func (s *Server) GetApps() ([]string, error) {
	apps := []string{}
	if err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucketName))
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