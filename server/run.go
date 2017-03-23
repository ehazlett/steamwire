package server

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *Server) Run() error {
	globalMux := http.NewServeMux()
	r, err := s.router()
	if err != nil {
		return err
	}
	globalMux.Handle("/", r)

	srv := &http.Server{
		Addr:    s.config.ListenAddr,
		Handler: globalMux,
	}

	logrus.WithFields(logrus.Fields{
		"addr":           s.config.ListenAddr,
		"updateInterval": s.config.UpdateInterval,
	}).Info("api started")

	// start ticker
	t := time.NewTicker(s.config.UpdateInterval)
	go func() {
		for range t.C {
			logrus.WithFields(logrus.Fields{
				"date": time.Now(),
			}).Info("updating news")
			s.Sync()
		}
	}()

	// update handler
	go func() {
		for {
			item := <-s.updateChan
			logrus.Debugf("update: %+v", item)
			// TODO: send to discord
			if err := s.sendToDiscord(item); err != nil {
				logrus.Errorf("error sending to discord: %s", err)
			}
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
