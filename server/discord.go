package server

import (
	"bytes"
	"html/template"

	"github.com/sirupsen/logrus"
)

const (
	msgTmpl = `**{{.Title}}** by {{.Author}}

{{.Contents}}

[Read more]({{.URL}})

[Application Page](http://store.steampowered.com/app/{{.AppID}}/)
`
)

func generateMessage(item *NewsItem) (string, error) {
	t := template.Must(template.New("message").Parse(msgTmpl))
	buf := &bytes.Buffer{}

	if err := t.Execute(buf, item); err != nil {
		return "", err
	}
	msg := buf.String()
	return msg, nil
}

func (s *Server) sendToDiscord(item *NewsItem) error {
	msg, err := generateMessage(item)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"message":        msg,
		"discordChannel": s.config.DiscordChannelID,
	}).Debug("sending to discord")

	if _, err := s.discord.ChannelMessageSend(s.config.DiscordChannelID, msg); err != nil {
		return err
	}

	return nil
}
