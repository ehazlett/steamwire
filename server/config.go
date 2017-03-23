package server

import "time"

// Config is the configuration for the server
type Config struct {
	ListenAddr       string
	DBPath           string
	UpdateInterval   time.Duration
	DiscordToken     string
	DiscordChannelID string
}
