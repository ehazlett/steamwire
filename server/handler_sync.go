package server

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (s *Server) processSyncMessage(author *discordgo.User, terms []string) (*handlerResponse, error) {
	logrus.WithFields(logrus.Fields{
		"author": author.Username,
	}).Debug("handling sync message")

	if err := s.sync(); err != nil {
		return &handlerResponse{
			Title:   errTitle,
			Content: errMessage,
		}, err
	}

	logrus.WithFields(logrus.Fields{
		"author": author.Username,
		"date":   time.Now(),
	}).Info("synchronized apps")
	return nil, nil
}
