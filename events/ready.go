package events

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
)

func Ready(bot *discordgo.Session, ready *discordgo.Ready) {
	transaction := sentry.StartTransaction(context.WithValue(context.Background(), clients.DiscordSession{}, bot), "ready")
	defer transaction.Finish()
	slog.Info(fmt.Sprintf("logged in as %s", ready.User.String()))

	clients.PublishNotification(transaction, fmt.Sprintf("logged in as %s", ready.User.String()))

	err := bot.UpdateCustomStatus(fmt.Sprintf("assisting %d classes", len(bot.State.Guilds)))
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to update status").Error())
	}
}
