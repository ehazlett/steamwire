package server

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/bwmarrin/discordgo"
)

const (
	dbBucketName = "apps"
)

type Server struct {
	config     *ServerConfig
	db         *bolt.DB
	updateChan chan (*NewsItem)
	discord    *discordgo.Session
}

func NewServer(cfg *ServerConfig) (*Server, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return nil, err
	}

	// ensure bucket is created
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(dbBucketName)); err != nil {
			return fmt.Errorf("error creating bucket: %s", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	ch := make(chan *NewsItem)

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
		db:         db,
		updateChan: ch,
		discord:    discord,
	}, nil
}
