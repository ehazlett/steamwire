package server

import (
	"sync"

	"github.com/sirupsen/logrus"
)

func (s *Server) Sync() error {
	apps, err := s.GetApps()
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	for _, app := range apps {
		wg.Add(1)
		go s.updateNewsForApp(app, wg, s.updateChan)
	}

	wg.Wait()
	return nil
}

func (s *Server) updateNewsForApp(appID string, wg *sync.WaitGroup, ch chan (*NewsItem)) {
	defer wg.Done()

	appNews, err := s.GetNews(appID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"appID": appID,
		}).Errorf("unable to update news for app: %s", err)
		return
	}

	if len(appNews.NewsItems) == 0 {
		logrus.WithFields(logrus.Fields{
			"appID": appID,
		}).Warnf("no news items in update")
		return
	}

	current, err := s.GetAppIndex(appID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"appID": appID,
		}).Errorf("error getting current app index: %s", err)
		return
	}

	item := appNews.NewsItems[0]
	if item.Gid == current {
		logrus.WithFields(logrus.Fields{
			"appID":   appID,
			"gid":     item.Gid,
			"current": current,
		}).Debug("news for app is current")
		return
	}
	if err := s.UpdateAppIndex(appID, item.Gid); err != nil {
		logrus.WithFields(logrus.Fields{
			"appID": appID,
			"gid":   item.Gid,
		}).Errorf("error updating news index: %s", err)
		return
	}

	// update
	ch <- item
}
