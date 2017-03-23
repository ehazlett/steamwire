package server

import "time"

type ServerConfig struct {
	ListenAddr       string
	DBPath           string
	UpdateInterval   time.Duration
	DiscordToken     string
	DiscordChannelID string
}
