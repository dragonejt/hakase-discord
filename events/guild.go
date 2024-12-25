package events

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/getsentry/sentry-go"
)

func GuildCreate(session *discordgo.Session, guildCreate *discordgo.GuildCreate) {
	sentry.StartTransaction(context.TODO(), "guildCreate")
	slog.Info(fmt.Sprintf("added to guild: %s (%s)", guildCreate.Guild.Name, guildCreate.Guild.ID))
	course := clients.Course{
		Course_id: guildCreate.Guild.ID,
	}
	err := clients.CreateCourse(course)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create course: %s", err))
	}

	err = session.UpdateCustomStatus(fmt.Sprintf("assisting %d classes", len(session.State.Guilds)))
	if err != nil {
		slog.Error(fmt.Sprintf("failed to update status: %s", err))
	}
}

func GuildDelete(session *discordgo.Session, guildDelete *discordgo.GuildDelete) {
	sentry.StartTransaction(context.TODO(), "guildDelete")
	slog.Info(fmt.Sprintf("removed from guild: %s (%s)", guildDelete.Guild.Name, guildDelete.Guild.ID))
	err := clients.DeleteCourse(guildDelete.Guild.ID)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to delete course: %s", err))
	}

	err = session.UpdateCustomStatus(fmt.Sprintf("assisting %d classes", len(session.State.Guilds)))
	if err != nil {
		slog.Error(fmt.Sprintf("failed to update status: %s", err))
	}
}
