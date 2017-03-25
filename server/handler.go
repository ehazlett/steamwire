package server

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ehazlett/steamwire/types"
	"github.com/sirupsen/logrus"
)

const (
	errTitle   = "oh noes!"
	errMessage = "something went wrong :scream_cat:"
)

type handlerResponse struct {
	Title   string
	Content string
	URL     string
}

func (s *Server) messageCreateHandler(session *discordgo.Session, m *discordgo.MessageCreate) {
	if !s.isBotMessage(m) {
		return
	}

	msgType, terms, err := s.getMessageData(m.Content)
	if err != nil {
		if err := s.sendDiscordMessage(errTitle, err.Error(), ""); err != nil {
			logrus.Errorf("error sending discord message: %s", err)
		}
		return
	}

	var f func(*discordgo.User, []string) (*handlerResponse, error)
	switch msgType {
	case types.MessageTypeAdd:
		f = s.processAddMessage
	case types.MessageTypeDelete:
		f = s.processDeleteMessage
	case types.MessageTypeSearch:
		f = s.processSearchMessage
	case types.MessageTypeList:
		f = s.processListMessage
	case types.MessageTypeSync:
		f = s.processSyncMessage
	case types.MessageTypeAnnounce:
		f = s.processAnnounceMessage
	case types.MessageTypeHelp:
		f = s.processHelpMessage
	}

	resp, err := f(m.Author, terms)
	if err != nil {
		logrus.Error(err)
	}
	if resp != nil {
		if err := s.sendDiscordMessage(resp.Title, resp.Content, resp.URL); err != nil {
			logrus.Errorf("error sending discord message: %s", err)
		}
	}
}

func (s *Server) isBotMessage(m *discordgo.MessageCreate) bool {
	for _, u := range m.Mentions {
		if u.Username == s.discordUser.Username {
			return true
		}
	}

	return false
}

// getMessageData parses the message content and returns the `types.MessageType`
// along with the query terms
func (s *Server) getMessageData(content string) (types.MessageType, []string, error) {
	fields := strings.Fields(content)
	if len(fields) == 1 {
		return types.MessageTypeUnknown, nil, fmt.Errorf("please specify an action")
	}

	term := fields[1]
	query := []string{}
	// default to all terms
	if len(fields) >= 3 {
		query = fields[2:]
	}

	switch strings.ToLower(term) {
	case string(types.MessageTypeAdd):
		return types.MessageTypeAdd, query, nil
	case string(types.MessageTypeDelete):
		return types.MessageTypeDelete, query, nil
	case string(types.MessageTypeList):
		return types.MessageTypeList, query, nil
	case string(types.MessageTypeSearch):
		return types.MessageTypeSearch, query, nil
	case string(types.MessageTypeSync):
		return types.MessageTypeSync, query, nil
	case string(types.MessageTypeHelp), "halp", "?", "halp!":
		return types.MessageTypeHelp, query, nil
	}

	return types.MessageTypeUnknown, query, fmt.Errorf("sorry I do not understand the action requested. :smiling_imp: ")
}
