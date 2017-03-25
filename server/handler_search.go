package server

import (
	"bytes"
	"html/template"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (s *Server) processSearchMessage(author *discordgo.User, terms []string) (*handlerResponse, error) {
	logrus.WithFields(logrus.Fields{
		"author": author.Username,
	}).Debug("handling search message")
	// ensure app list is updated
	if err := s.updateAppList(false); err != nil {
		return &handlerResponse{
			Title:   errTitle,
			Content: errMessage,
		}, err
	}
	query := strings.Join(terms, " ")
	apps, err := s.ds.FindApp(query)
	if err != nil {
		return &handlerResponse{
			Title:   errTitle,
			Content: errMessage,
		}, err
	}
	if len(apps) == 0 {
		return &handlerResponse{
			Title:   errTitle,
			Content: "sorry I was unable to find any applications matching that search :slight_frown:",
		}, nil
	}
	if len(apps) > 30 {
		return &handlerResponse{
			Title:   errTitle,
			Content: "sorry that returns too many results and I do not want to spam the channel.  please refine your search. :slight_smile:",
		}, nil
	}
	tmpl := `{{range .}}[{{.Name}}](https://store.steampowered.com/app/{{.AppID}}) ({{.AppID}})
{{end}}
	`
	t := template.Must(template.New("message").Parse(tmpl))
	buf := &bytes.Buffer{}

	if err := t.Execute(buf, apps); err != nil {
		return &handlerResponse{
			Title:   errTitle,
			Content: errMessage,
		}, err
	}
	msg := buf.String()
	logrus.WithFields(logrus.Fields{
		"author":    author.Username,
		"query":     query,
		"numOfApps": len(apps),
		"date":      time.Now(),
	}).Info("application search")
	return &handlerResponse{
		Title:   "Search Results",
		Content: msg,
	}, nil
}
