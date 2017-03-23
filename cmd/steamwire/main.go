package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/pkg/homedir"
	"github.com/ehazlett/steamwire/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func getDefaultDBPath() string {
	h := homedir.Get()
	return filepath.Join(h, ".steamwire.db")
}

func main() {
	app := cli.NewApp()
	app.Name = version.Name()
	app.Usage = version.Description()
	app.Version = version.Version()
	app.Author = "@ehazlett"
	app.Email = ""
	app.Before = func(c *cli.Context) error {
		// enable debug
		if c.GlobalBool("debug") {
			log.SetLevel(log.DebugLevel)
			log.Debug("debug enabled")
		}

		return nil
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "enable debug",
		},
		cli.StringFlag{
			Name:  "db-path, d",
			Usage: "Path to database",
			Value: getDefaultDBPath(),
		},
		cli.StringFlag{
			Name:  "listen-addr, l",
			Usage: "Listen address",
			Value: ":8080",
		},
		cli.DurationFlag{
			Name:  "update-interval, i",
			Usage: "Update interval",
			Value: time.Hour * 1,
		},
		cli.StringFlag{
			Name:   "discord-token, t",
			Usage:  "Discord bot token",
			Value:  "",
			EnvVar: "DISCORD_TOKEN",
		},
		cli.StringFlag{
			Name:   "discord-channel-id, c",
			Usage:  "Discord channel ID",
			Value:  "",
			EnvVar: "DISCORD_CHANNEL_ID",
		},
	}
	app.Action = runAction

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
