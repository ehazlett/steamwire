package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ehazlett/steamwire/types"
	"github.com/sirupsen/logrus"
)

const (
	searchURL = "http://api.steampowered.com/ISteamApps/GetAppList/v0001/"
)

// updateAppList will only update the local app list once an hour unless forced
func (s *Server) updateAppList(force bool) error {
	// check last refresh date
	if !force {
		lastUpdated, err := s.ds.GetAppListLastUpdated()
		if err != nil {
			return err
		}

		// check if less than an hour; if so return
		diff := time.Now().Sub(lastUpdated)
		if diff < time.Duration(1*time.Hour) {
			logrus.WithFields(logrus.Fields{
				"lastUpdated": lastUpdated,
			}).Debug("skipping app update; was last updated within the hour")
			return nil
		}
	}
	logrus.WithFields(logrus.Fields{
		"time": time.Now(),
	}).Debug("starting app list update")
	resp, err := http.Get(searchURL)
	if err != nil {
		return err
	}
	var appList *types.List
	if err := json.NewDecoder(resp.Body).Decode(&appList); err != nil {
		return err
	}

	apps := appList.AppList.Apps.Info
	if err := s.ds.UpdateAppList(apps); err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"numOfApps": len(apps),
		"time":      time.Now(),
	}).Info("updated app list")

	return nil
}
