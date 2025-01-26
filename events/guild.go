package events

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/getsentry/sentry-go"
)

func GuildCreate(bot *discordgo.Session, guildCreate *discordgo.GuildCreate) {
	transaction := sentry.StartTransaction(context.WithValue(context.Background(), clients.DiscordSession{}, bot), "guildCreate")
	defer transaction.Finish()
	slog.Info(fmt.Sprintf("added to guild: %s (%s)", guildCreate.Guild.Name, guildCreate.Guild.ID))

	course := clients.Course{
		CourseID: guildCreate.Guild.ID,
	}
	err := clients.CreateCourse(transaction, course)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create course: %s", err))
	}

	err = bot.UpdateCustomStatus(fmt.Sprintf("assisting %d classes", len(bot.State.Guilds)))
	if err != nil {
		slog.Error(fmt.Sprintf("failed to update status: %s", err))
	}
}

func GuildDelete(bot *discordgo.Session, guildDelete *discordgo.GuildDelete) {
	transaction := sentry.StartTransaction(context.WithValue(context.Background(), clients.DiscordSession{}, bot), "guildDelete")
	defer transaction.Finish()
	slog.Info(fmt.Sprintf("removed from guild: %s (%s)", guildDelete.Guild.Name, guildDelete.Guild.ID))

	err := clients.DeleteCourse(transaction, guildDelete.Guild.ID)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to delete course: %s", err))
	}

	err = bot.UpdateCustomStatus(fmt.Sprintf("assisting %d classes", len(bot.State.Guilds)))
	if err != nil {
		slog.Error(fmt.Sprintf("failed to update status: %s", err))
	}
}
