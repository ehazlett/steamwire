package server

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Run starts the server
func (s *Server) Run() error {
	// start ticker
	t := time.NewTicker(s.config.UpdateInterval)
	go func() {
		for range t.C {
			logrus.WithFields(logrus.Fields{
				"date": time.Now(),
			}).Info("updating news")
			s.sync()
		}
	}()

	// update handler
	go func() {
		for {
			item := <-s.updateChan
			logrus.WithFields(logrus.Fields{
				"appID": item.AppID,
				"gid":   item.Gid,
				"date":  time.Now(),
			}).Debugf("app update")
			// send to discord
			if err := s.sendToDiscord(item); err != nil {
				logrus.Errorf("error sending to discord: %s", err)
			}
		}
	}()

	// ensure connected
	if err := s.ensureConnectionToDiscord(); err != nil {
		return err
	}
	user, err := s.discord.User("@me")
	if err != nil {
		return err
	}
	s.discordUser = user

	// add handler
	s.discord.AddHandler(s.messageCreateHandler)

	return nil
}
