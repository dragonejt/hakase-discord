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
	sentry.StartTransaction(context.Background(), "guildCreate")
	slog.Info(fmt.Sprintf("Added to Server: %s (%s)", guildCreate.Guild.Name, guildCreate.Guild.ID))
	course := clients.Course{
		CourseID: guildCreate.Guild.ID,
	}
	err := clients.CreateCourse(course)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to Create Course: %s", err))
	}

	err = bot.UpdateCustomStatus(fmt.Sprintf("Assisting %d Classes", len(bot.State.Guilds)))
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to Update Status: %s", err))
	}
}

func GuildDelete(bot *discordgo.Session, guildDelete *discordgo.GuildDelete) {
	sentry.StartTransaction(context.Background(), "guildDelete")
	slog.Info(fmt.Sprintf("Removed from Server: %s (%s)", guildDelete.Guild.Name, guildDelete.Guild.ID))
	err := clients.DeleteCourse(guildDelete.Guild.ID)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to Delete Course: %s", err))
	}

	err = bot.UpdateCustomStatus(fmt.Sprintf("Assisting %d Classes", len(bot.State.Guilds)))
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to Update Status: %s", err))
	}
}
