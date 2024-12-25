package events

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/notifications"
	"github.com/getsentry/sentry-go"
)

func Ready(session *discordgo.Session, ready *discordgo.Ready) {
	sentry.StartTransaction(context.TODO(), "ready")
	slog.Info(fmt.Sprintf("logged in as %s", ready.User.String()))
	notifications.PublishNotification(fmt.Sprintf("logged in as %s", ready.User.String()))

	err := session.UpdateCustomStatus(fmt.Sprintf("assisting %d classes", len(session.State.Guilds)))
	if err != nil {
		slog.Error(fmt.Sprintf("failed to update status: %s", err))
	}
}
