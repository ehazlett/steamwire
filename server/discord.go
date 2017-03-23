package server

import "fmt"

func (s *Server) sendToDiscord(item *NewsItem) error {
	msg := fmt.Sprintf("**%s** by %s\n\n%s", item.Title, item.Author, item.Contents)
	if _, err := s.discord.ChannelMessageSend(s.config.DiscordChannelID, msg); err != nil {
		return err
	}

	return nil
}
