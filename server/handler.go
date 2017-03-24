package server

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ehazlett/steamwire/types"
	"github.com/sirupsen/logrus"
)

const (
	errMessage = "oh no! something went wrong :("
)

func (s *Server) messageCreateHandler(session *discordgo.Session, m *discordgo.MessageCreate) {
	if !s.isBotMessage(m) {
		return
	}

	author := m.Author.Username

	msgType, query, err := s.getMessageData(m.Content)
	if err != nil {
		if err := s.sendDiscordMessage(err.Error()); err != nil {
			logrus.Errorf("error sending discord message: %s", err)
		}
		return
	}

	var f func(string, string) (string, error)
	switch msgType {
	case types.MessageTypeAdd:
		f = s.processAddMessage
	case types.MessageTypeDelete:
		f = s.processDeleteMessage
	case types.MessageTypeSearch:
		f = s.processSearchMessage
	case types.MessageTypeList:
		f = s.processListMessage
	case types.MessageTypeHelp:
		f = s.processHelpMessage
	}

	resp, err := f(author, query)
	if err != nil {
		logrus.Error(err)
	}
	if err := s.sendDiscordMessage(resp); err != nil {
		logrus.Errorf("error sending discord message: %s", err)
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

func (s *Server) getMessageData(content string) (types.MessageType, string, error) {
	fields := strings.Fields(content)
	if len(fields) == 1 {
		return types.MessageTypeUnknown, "", fmt.Errorf("please specify an action")
	}

	term := fields[1]
	query := ""

	switch strings.ToLower(term) {
	case string(types.MessageTypeAdd):
		if len(fields) == 3 {
			query = fields[2]
		}
		return types.MessageTypeAdd, query, nil
	case string(types.MessageTypeDelete):
		if len(fields) == 3 {
			query = fields[2]
		}
		return types.MessageTypeDelete, query, nil
	case string(types.MessageTypeList):
		return types.MessageTypeList, query, nil
	case string(types.MessageTypeSearch):
		if len(fields) >= 3 {
			query = strings.Join(fields[2:], " ")
		}
		return types.MessageTypeSearch, query, nil
	case string(types.MessageTypeHelp), "halp", "?", "halp!":
		return types.MessageTypeHelp, "", nil
	}

	return types.MessageTypeUnknown, "", fmt.Errorf("sorry I do not understand the action requested. :smiling_imp: ")
}

func (s *Server) processAddMessage(author, appID string) (string, error) {
	logrus.WithFields(logrus.Fields{
		"author": author,
	}).Debug("handling add message")

	info, err := s.ds.GetAppInfo(appID)
	if err != nil {
		return errMessage, err
	}
	if info == nil {
		return fmt.Sprintf("sorry that does not appear to be a valid app id"), nil
	}

	if err := s.ds.AddApp(appID); err != nil {
		return fmt.Sprintf("sorry unable to add application: %s", err), err
	}

	logrus.WithFields(logrus.Fields{
		"appID":  appID,
		"name":   info.Name,
		"author": author,
	}).Info("added application")

	return fmt.Sprintf("%s (%s) added!", info.Name, appID), nil
}

func (s *Server) processDeleteMessage(author, appID string) (string, error) {
	logrus.WithFields(logrus.Fields{
		"author": author,
	}).Debug("handling delete message")

	info, err := s.ds.GetAppInfo(appID)
	if err != nil {
		return errMessage, err
	}
	if info == nil {
		return fmt.Sprintf("sorry that does not appear to be a valid app id."), nil
	}

	// check if being watched
	apps, err := s.ds.GetApps()
	if err != nil {
		return errMessage, err
	}

	found := false
	for _, app := range apps {
		if app == appID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Sprintf("sorry that application is not being monitored"), nil
	}

	if err := s.ds.DeleteApp(appID); err != nil {
		return fmt.Sprintf("sorry unable to delete application: %s", err), err
	}

	logrus.WithFields(logrus.Fields{
		"appID":  appID,
		"name":   info.Name,
		"author": author,
	}).Info("deleted application")

	return fmt.Sprintf("%s (%s) removed!", info.Name, appID), nil
}

func (s *Server) processSearchMessage(author, query string) (string, error) {
	logrus.WithFields(logrus.Fields{
		"author": author,
	}).Debug("handling search message")
	apps, err := s.ds.FindApp(query)
	if err != nil {
		return errMessage, err
	}
	if len(apps) == 0 {
		return fmt.Sprintf("sorry I was unable to find any applications matching that search :slight_frown:"), nil
	}
	if len(apps) > 30 {
		return fmt.Sprintf("sorry that returns too many results and I do not want to spam the channel.  please refine your search. :slight_smile:"), nil
	}
	tmpl := `

**Results**

{{range .}}{{.Name}} ({{.AppID}})
{{end}}
	`
	t := template.Must(template.New("message").Parse(tmpl))
	buf := &bytes.Buffer{}

	if err := t.Execute(buf, apps); err != nil {
		return "", err
	}
	msg := buf.String()
	logrus.WithFields(logrus.Fields{
		"author":    author,
		"query":     query,
		"numOfApps": len(apps),
	}).Info("application search")
	return msg, nil
}

func (s *Server) processListMessage(author, query string) (string, error) {
	logrus.WithFields(logrus.Fields{
		"author": author,
	}).Debug("handling list message")

	apps, err := s.ds.GetApps()
	if err != nil {
		return errMessage, err
	}

	info := []*types.AppInfo{}
	for _, app := range apps {
		i, err := s.ds.GetAppInfo(app)
		if err != nil {
			return errMessage, err
		}

		info = append(info, i)
	}

	tmpl := `

**Current Applications**
{{range .}}{{.Name}} ({{.AppID}})
{{end}}
	`
	t := template.Must(template.New("message").Parse(tmpl))
	buf := &bytes.Buffer{}

	if err := t.Execute(buf, info); err != nil {
		return "", err
	}
	msg := buf.String()
	return msg, nil
}

func (s *Server) processHelpMessage(author, query string) (string, error) {
	logrus.WithFields(logrus.Fields{
		"author": author,
	}).Debug("handling help message")
	msg := `
Sure I can help!  Here are a list of commands that I currently process:
` + "```" + `

list		    List currently monitored applications
add <id>	    Add a new application to monitored list
delete <id>	    Delete a monitored application
search <query>	    Search for an application
help		    (this message) :)
` + "```"

	return msg, nil
}
