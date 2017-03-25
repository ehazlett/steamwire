package server

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (s *Server) processDeleteMessage(author *discordgo.User, terms []string) (*handlerResponse, error) {
	logrus.WithFields(logrus.Fields{
		"author": author.Username,
	}).Debug("handling delete message")
	if len(terms) == 0 {
		return &handlerResponse{
			Title:   errTitle,
			Content: "sorry you must specify an application id",
		}, nil
	}
	appID := terms[0]

	info, err := s.ds.GetAppInfo(appID)
	if err != nil {
		return &handlerResponse{
			Title:   errTitle,
			Content: errMessage,
		}, err
	}
	if info == nil {
		return &handlerResponse{
			Title:   errTitle,
			Content: "sorry that does not appear to be a valid app id",
		}, nil
	}

	// check if being watched
	apps, err := s.ds.GetApps()
	if err != nil {
		return &handlerResponse{
			Title:   errTitle,
			Content: errMessage,
		}, err
	}

	found := false
	for _, app := range apps {
		if app == appID {
			found = true
			break
		}
	}

	if !found {
		return &handlerResponse{
			Title:   errTitle,
			Content: "sorry that application is not being monitored",
		}, nil
	}

	if err := s.ds.DeleteApp(appID); err != nil {
		return &handlerResponse{
			Title:   errTitle,
			Content: fmt.Sprintf("sorry unable to delete application: %s", err),
		}, nil
	}

	logrus.WithFields(logrus.Fields{
		"appID":  appID,
		"name":   info.Name,
		"author": author.Username,
	}).Info("deleted application")

	resp := &handlerResponse{
		Content: fmt.Sprintf("<@!%s> %s (%s) deleted!", author.ID, info.Name, appID),
	}

	return resp, nil
}
