package server

import (
	"bytes"
	"html/template"

	"github.com/ehazlett/steamwire/types"
	"github.com/sirupsen/logrus"
)

const (
	msgTmpl = `**{{.Title}}** by {{.Author}}

{{.Contents}}

Read more: {{.URL}}

Application Page: http://store.steampowered.com/app/{{.AppID}}/
`
)

func generateMessage(item *types.NewsItem) (string, error) {
	t := template.Must(template.New("message").Parse(msgTmpl))
	buf := &bytes.Buffer{}

	if err := t.Execute(buf, item); err != nil {
		return "", err
	}
	msg := buf.String()
	return msg, nil
}

func (s *Server) sendToDiscord(item *types.NewsItem) error {
	if s.discord == nil {
		logrus.Warnf("discord is not configured; skipping send")
		return nil
	}

	msg, err := generateMessage(item)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"appID":          item.AppID,
		"title":          item.Title,
		"gid":            item.Gid,
		"discordChannel": s.config.DiscordChannelID,
	}).Debug("sending to discord")

	if err := s.sendDiscordMessage(msg); err != nil {
		return err
	}

	return nil
}

func (s *Server) sendDiscordMessage(content string) error {
	if _, err := s.discord.ChannelMessageSend(s.config.DiscordChannelID, content); err != nil {
		return err
	}

	return nil
}
