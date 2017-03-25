package server

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (s *Server) processAddMessage(author *discordgo.User, terms []string) (*handlerResponse, error) {
	logrus.WithFields(logrus.Fields{
		"author": author.Username,
	}).Debug("handling add message")
	// make sure an id is specified
	if len(terms) == 0 {
		return &handlerResponse{
			Title:   errTitle,
			Content: "sorry you must specify an application id",
		}, nil
	}

	appID := terms[0]
	// ensure app list is updated
	if err := s.updateAppList(false); err != nil {
		return &handlerResponse{
			Title:   errTitle,
			Content: errMessage,
		}, err
	}

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
			Content: fmt.Sprintf("sorry that does not appear to be a valid app id"),
		}, nil
	}

	if err := s.ds.AddApp(appID); err != nil {
		return &handlerResponse{
			Title:   errTitle,
			Content: fmt.Sprintf("sorry unable to add application: %s", err),
		}, err
	}

	resp := &handlerResponse{
		Content: fmt.Sprintf("<@!%s> %s (%s) added!", author.ID, info.Name, appID),
	}

	// attempt to sync; log errors to app only -- don't send to user
	appNews, err := s.getNews(appID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"appID": appID,
			"name":  info.Name,
		}).Error("error getting news for app")
		return resp, err
	}
	items := appNews.NewsItems
	if len(items) > 0 {
		item := items[0]
		updated, err := s.updateNewsForApp(appID, item)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"appID": appID,
				"name":  info.Name,
			}).Error("error updating news for app")
			return resp, err
		}

		if updated {
			if err := s.sendToDiscord(item); err != nil {
				logrus.WithFields(logrus.Fields{
					"appID": appID,
					"name":  info.Name,
				}).Error("error sending update for app")
				return resp, err
			}
		}
	}

	logrus.WithFields(logrus.Fields{
		"appID":  appID,
		"name":   info.Name,
		"author": author.Username,
	}).Info("added application")

	return resp, nil
}
