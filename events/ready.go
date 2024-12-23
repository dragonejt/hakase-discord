package events

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/notifications"
)

func Ready(session *discordgo.Session, ready *discordgo.Ready) {
	slog.Info(fmt.Sprintf("logged in as %s", ready.User.String()))
	notifications.PublishNotification(fmt.Sprintf("logged in as %s", ready.User.String()))
	session.UpdateCustomStatus(fmt.Sprintf("assisting %d classes", len(session.State.Guilds)))
}
