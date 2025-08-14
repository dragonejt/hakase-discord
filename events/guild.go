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

func GuildCreate(bot *discordgo.Session, guildCreate *discordgo.GuildCreate, hakaseClient clients.HakaseClient) {
	transaction := sentry.StartTransaction(context.WithValue(context.Background(), clients.DiscordSession{}, bot), "guildCreate")
	defer transaction.Finish()
	slog.Info(fmt.Sprintf("added to guild: %s (%s)", guildCreate.Guild.Name, guildCreate.Guild.ID))

	course := clients.Course{
		CourseID: guildCreate.Guild.ID,
	}
	err := hakaseClient.CreateCourse(transaction, course)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to create course").Error())
	}

	err = bot.UpdateCustomStatus(fmt.Sprintf("assisting %d classes", len(bot.State.Guilds)))
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to update status").Error())
	}
}

func GuildDelete(bot *discordgo.Session, guildDelete *discordgo.GuildDelete, hakaseClient clients.HakaseClient) {
	transaction := sentry.StartTransaction(context.WithValue(context.Background(), clients.DiscordSession{}, bot), "guildDelete")
	defer transaction.Finish()
	slog.Info(fmt.Sprintf("removed from guild: %s (%s)", guildDelete.Guild.Name, guildDelete.Guild.ID))

	err := hakaseClient.DeleteCourse(transaction, guildDelete.Guild.ID)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to delete course").Error())
	}

	err = bot.UpdateCustomStatus(fmt.Sprintf("assisting %d classes", len(bot.State.Guilds)))
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to update status").Error())
	}
}
