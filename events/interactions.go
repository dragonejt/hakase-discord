package events

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/interactions"
)

func InteractionCreate(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	switch interactionCreate.Type {
	case discordgo.InteractionApplicationCommand:
		switch interactionCreate.ApplicationCommandData().Name {
		case "assignments":
			interactions.SlashAssignments(bot, interactionCreate)
		case "hakase":
			interactions.SlashHakase(bot, interactionCreate)
		default:
			slog.Error(fmt.Sprintf("unknown command: %s", interactionCreate.ApplicationCommandData().Name))
		}
	case discordgo.InteractionMessageComponent:
		switch interactionCreate.MessageComponentData().CustomID {
		case "addAssignmentAction":
			interactions.AddAssignment(bot, interactionCreate)
		default:
			slog.Error(fmt.Sprintf("unknown message component: %s", interactionCreate.MessageComponentData().CustomID))
		}
	case discordgo.InteractionModalSubmit:
		switch interactionCreate.ModalSubmitData().CustomID {
		case "addAssignment":
			interactions.AddAssignmentSubmit(bot, interactionCreate)
		default:
			slog.Error(fmt.Sprintf("unknown modal submit: %s", interactionCreate.ModalSubmitData().CustomID))
		}
	default:
		slog.Error(fmt.Sprintf("unknown interaction type: %d", interactionCreate.Type))

	}
}
