package server

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (s *Server) processHelpMessage(author *discordgo.User, terms []string) (*handlerResponse, error) {
	logrus.WithFields(logrus.Fields{
		"author": author.Username,
	}).Debug("handling help message")
	msg := `
Sure I can help!  Here are a list of commands that I currently understand:
` + "```" + `

list                List currently monitored applications
add <id>            Add a new application to monitored list
delete <id>         Delete a monitored application
search <query>      Search for an application
help                (this message) :)
` + "```"
	logrus.WithFields(logrus.Fields{
		"author": author.Username,
		"date":   time.Now(),
	}).Info("requested help")

	return &handlerResponse{
		Content: msg,
	}, nil
}
