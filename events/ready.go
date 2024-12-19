package events

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/notifications"
)

func Ready(session *discordgo.Session, ready *discordgo.Ready) {
	slog.Info("logged in as " + ready.User.String())
	notifications.PublishNotification("logged in as " + ready.User.String())
}
