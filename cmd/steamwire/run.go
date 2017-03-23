package main

import (
	"github.com/ehazlett/steamwire/server"
	"github.com/ehazlett/steamwire/version"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	registrationLink = "https://discordapp.com/oauth2/authorize?client_id=294309210890436610&scope=bot&permissions=19456"
)

func runAction(c *cli.Context) error {
	logrus.Info(version.FullVersion())

	discordToken := c.GlobalString("discord-token")
	discordChannel := c.GlobalString("discord-channel-id")

	if discordToken == "" || discordChannel == "" {
		logrus.Info("  ************************************************************************************")
		logrus.Infof("  Visit %s to authorize to your Discord guild", registrationLink)
		logrus.Infof("  Please be sure to set the channel ID from the channel you would like to use as well")
		logrus.Info("  ************************************************************************************")
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
