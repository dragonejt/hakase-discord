package events

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/commands"
)

func InteractionCreate(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	switch interactionCreate.Type {
	case discordgo.InteractionApplicationCommand:
		switch interactionCreate.ApplicationCommandData().Name {
		case "assignments":
			commands.SlashAssignments(bot, interactionCreate)
		case "hakase":
			commands.SlashHakase(bot, interactionCreate)
		default:
			slog.Error(fmt.Sprintf("unknown command: %s", interactionCreate.ApplicationCommandData().Name))
		}
	case discordgo.InteractionMessageComponent:
		slog.Error("message component interaction not implemented")
	default:
		slog.Error(fmt.Sprintf("unknown interaction type: %d", interactionCreate.Type))

	}
}
