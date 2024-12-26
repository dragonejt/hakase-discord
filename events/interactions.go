package events

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/commands"
)

func InteractionCreate(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	switch interactionCreate.ApplicationCommandData().Name {
	case "assignments":
		commands.SlashAssignments(bot, interactionCreate)
	case "hakase":
		commands.SlashHakase(bot, interactionCreate)
	default:
		slog.Error(fmt.Sprintf("unknown command: %s", interactionCreate.ApplicationCommandData().Name))
	}
}
