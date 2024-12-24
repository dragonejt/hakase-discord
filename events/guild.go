package events

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
)

func GuildCreate(session *discordgo.Session, guildCreate *discordgo.GuildCreate) {
	slog.Info(fmt.Sprintf("added to guild: %s (%s)", guildCreate.Guild.Name, guildCreate.Guild.ID))
	course := clients.Course{
		Course_id: guildCreate.Guild.ID,
	}
	err := clients.CreateCourse(course)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create course: %s", err))
		return
	}

	err = session.UpdateCustomStatus(fmt.Sprintf("assisting %d classes", len(session.State.Guilds)))
	if err != nil {
		slog.Error(fmt.Sprintf("failed to update status: %s", err))
		return
	}
}

func GuildDelete(session *discordgo.Session, guildDelete *discordgo.GuildDelete) {
	slog.Info(fmt.Sprintf("removed from guild: %s (%s)", guildDelete.Guild.Name, guildDelete.Guild.ID))
	err := clients.DeleteCourse(guildDelete.Guild.ID)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to delete course: %s", err))
		return
	}

	err = session.UpdateCustomStatus(fmt.Sprintf("assisting %d classes", len(session.State.Guilds)))
	if err != nil {
		slog.Error(fmt.Sprintf("failed to update status: %s", err))
		return
	}
}
