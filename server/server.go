package server

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/bwmarrin/discordgo"
	"github.com/ehazlett/steamwire/db"
	"github.com/ehazlett/steamwire/types"
)

// Server is the object for the core server
type Server struct {
	config      *Config
	ds          *db.DB
	updateChan  chan (*types.NewsItem)
	discord     *discordgo.Session
	discordUser *discordgo.User
}

// NewServer returns a new `Server`
func NewServer(cfg *Config) (*Server, error) {
	bdb, err := bolt.Open(cfg.DBPath, 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return nil, err
	}

	ds, err := db.NewDB(bdb)
	if err != nil {
		return nil, err
	}

	ch := make(chan *types.NewsItem)

	var discord *discordgo.Session
	if cfg.DiscordToken != "" {
		d, err := discordgo.New("Bot " + cfg.DiscordToken)
		if err != nil {
			return nil, err
		}

		if err := d.Open(); err != nil {
			return nil, err
		}

		discord = d
	}

	return &Server{
		config:     cfg,
		ds:         ds,
		updateChan: ch,
		discord:    discord,
	}, nil
}
