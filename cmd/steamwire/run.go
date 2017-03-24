package main

import (
	"fmt"
	"os"

	"github.com/ehazlett/steamwire/server"
	"github.com/ehazlett/steamwire/version"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func runAction(c *cli.Context) error {
	logrus.Info(version.FullVersion())

	discordToken := c.GlobalString("discord-token")
	discordChannel := c.GlobalString("discord-channel-id")

	if discordToken == "" || discordChannel == "" && os.Getenv("STEAMWIRE_DEBUG") == "" {
		help := `
    ------------------------------------------------------
    Please visit https://github.com/ehazlett/steamwire/blob/master/docs/install.md
    to get Steamwire configured with your Discord Guild.

    Once ready, specify the following options:

    --discord-token <your-discord-token>
    --discord-channel-id <your-channel-id>

    ------------------------------------------------------
`
		fmt.Println(help)
		return nil
	}

	cfg := &server.Config{
		ListenAddr:       c.GlobalString("listen-addr"),
		UpdateInterval:   c.GlobalDuration("update-interval"),
		DBPath:           c.GlobalString("db-path"),
		DiscordToken:     discordToken,
		DiscordChannelID: discordChannel,
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		return err
	}

	return srv.Run()
}
