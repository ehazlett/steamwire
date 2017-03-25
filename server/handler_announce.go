package server

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (s *Server) processAnnounceMessage(author *discordgo.User, terms []string) (*handlerResponse, error) {
	logrus.WithFields(logrus.Fields{
		"author": author.Username,
	}).Debug("handling announce message")

	if len(terms) == 0 {
		return &handlerResponse{
			Title:   errTitle,
			Content: "sorry you must specify an application id",
		}, nil
	}
	//appID := terms[0]

	//if err := s.sync(); err != nil {
	//	return "", err
	//}
	return nil, nil
}
