package server

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ehazlett/steamwire/types"
	"github.com/jaytaylor/html2text"
	"github.com/sirupsen/logrus"
)

const (
	msgTmpl = `
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
	// convert html to text
	txt, err := html2text.FromReader(buf)
	if err != nil {
		return "", err
	}
	return txt, nil
}

func (s *Server) ensureConnectionToDiscord() error {
	if _, err := s.discord.Gateway(); err != nil {
		// attempt to reconnect and try again
		if err := s.discord.Open(); err != nil {
			return err
		}

		g, err := s.discord.Gateway()
		if err != nil {
			return err
		}
		logrus.WithFields(logrus.Fields{
			"gateway": g,
			"date":    time.Now(),
		}).Debug("reconnected to discord")
		return nil
	}

	return nil
}

func (s *Server) sendToDiscord(item *types.NewsItem) error {
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

	if err := s.sendDiscordMessage(item.Title, msg, item.URL); err != nil {
		return err
	}

	return nil
}

func (s *Server) sendDiscordMessage(title string, content string, url string) error {
	if err := s.ensureConnectionToDiscord(); err != nil {
		return err
	}

	// check which type of message to send based upon params
	if title != "" || url != "" {
		// check for length and ellipsize as description has a max count
		// and will return an error
		fields := strings.Fields(content)
		if len(fields) > 100 {
			c := strings.Join(fields[:100], " ") + "..."
			content = c
			if url != "" {
				content = fmt.Sprintf("%s [read more](%s)", c, url)
			}
		}

		msg := &discordgo.MessageEmbed{
			Type:        "rich",
			Description: content,
		}
		if title != "" {
			msg.Title = title
		}
		if url != "" {
			msg.URL = url
		}

		if _, err := s.discord.ChannelMessageSendEmbed(s.config.DiscordChannelID, msg); err != nil {
			return err
		}
	} else {
		if _, err := s.discord.ChannelMessageSend(s.config.DiscordChannelID, content); err != nil {
			return err
		}
	}

	return nil
}
