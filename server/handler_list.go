package server

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/bwmarrin/discordgo"
	"github.com/ehazlett/steamwire/types"
	"github.com/sirupsen/logrus"
)

func (s *Server) processListMessage(author *discordgo.User, terms []string) (*handlerResponse, error) {
	logrus.WithFields(logrus.Fields{
		"author": author.Username,
	}).Debug("handling list message")

	apps, err := s.ds.GetApps()
	if err != nil {
		return &handlerResponse{
			Title:   errTitle,
			Content: errMessage,
		}, err
	}

	if len(apps) == 0 {
		return &handlerResponse{
			Content: fmt.Sprintf("<@!%s> There are no monitored applications.", author.ID),
		}, nil
	}

	info := []*types.AppInfo{}
	for _, app := range apps {
		i, err := s.ds.GetAppInfo(app)
		if err != nil {
			return &handlerResponse{
				Title:   errTitle,
				Content: errMessage,
			}, err
		}

		info = append(info, i)
	}

	tmpl := `{{range .}}[{{.Name}}](https://store.steampowered.com/app/{{.AppID}}/) ({{.AppID}})
{{end}}
	`
	t := template.Must(template.New("message").Parse(tmpl))
	buf := &bytes.Buffer{}

	if err := t.Execute(buf, info); err != nil {
		return &handlerResponse{
			Title:   errTitle,
			Content: errMessage,
		}, err
	}
	msg := buf.String()
	return &handlerResponse{
		Title:   "Applications",
		Content: msg,
	}, nil
}
