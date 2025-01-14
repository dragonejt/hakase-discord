package events

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/notifications"
	"github.com/getsentry/sentry-go"
)

func Ready(bot *discordgo.Session, ready *discordgo.Ready) {
	sentry.StartTransaction(context.Background(), "ready")
	slog.Info(fmt.Sprintf("Logged In as %s", ready.User.String()))
	notifications.PublishNotification(fmt.Sprintf("Logged In as %s", ready.User.String()))

	err := bot.UpdateCustomStatus(fmt.Sprintf("Assisting %d Classes", len(bot.State.Guilds)))
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to Update Status: %s", err))
	}
}
