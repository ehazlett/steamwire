package server

import (
	"fmt"
	"sync"

	"github.com/ehazlett/steamwire/types"
	"github.com/sirupsen/logrus"
)

// sync gets the latest news item for all applications in the database
func (s *Server) sync() error {
	apps, err := s.ds.GetApps()
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	for _, app := range apps {
		wg.Add(1)
		go s.syncNews(app, wg, s.updateChan)
	}

	wg.Wait()
	return nil
}

func (s *Server) syncNews(appID string, wg *sync.WaitGroup, ch chan (*types.NewsItem)) {
	defer wg.Done()

	appNews, err := s.getNews(appID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"appID": appID,
		}).Error("unable to update news for app")
		return
	}

	items := appNews.NewsItems
	if len(items) == 0 {
		logrus.WithFields(logrus.Fields{
			"appID": appID,
		}).Debugf("no news items in update")
		return
	}
	item := items[0]

	updated, err := s.updateNewsForApp(appID, item)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"appID": appID,
		}).Errorf("error syncing news for app")
		return
	}

	// update
	if updated {
		ch <- item
	}
}

// updateNewsForApp sets the latest update index for the app and returns
// whether or not it was updated
func (s *Server) updateNewsForApp(appID string, item *types.NewsItem) (bool, error) {
	current, err := s.ds.GetAppIndex(appID)
	if err != nil {
		return false, fmt.Errorf("error getting current app index: %s", err)
	}

	if item.Gid == current {
		logrus.WithFields(logrus.Fields{
			"appID":   appID,
			"gid":     item.Gid,
			"current": current,
		}).Debug("news for app is current")
		return false, nil
	}
	if err := s.ds.UpdateAppIndex(appID, item.Gid); err != nil {
		return false, fmt.Errorf("error updating news index: %s", err)
	}

	return true, nil
}
