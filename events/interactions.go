package events

import (
	"fmt"
	"log/slog"
	"strings"

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
		customID := interactionCreate.MessageComponentData().CustomID
		if strings.HasPrefix(customID, "addAssignmentAction") {
			interactions.AddAssignment(bot, interactionCreate)
		} else if strings.HasPrefix(customID, "updateAssignmentAction") {
			interactions.UpdateAssignment(bot, interactionCreate)
		} else if strings.HasPrefix(customID, "deleteAssignmentAction") {
			interactions.DeleteAssignment(bot, interactionCreate)
		} else {
			slog.Error(fmt.Sprintf("unknown message component action: %s", customID))
		}
	case discordgo.InteractionModalSubmit:
		customID := interactionCreate.ModalSubmitData().CustomID
		if strings.HasPrefix(customID, "addAssignment") {
			interactions.AddAssignmentSubmit(bot, interactionCreate)
		} else if strings.HasPrefix(customID, "updateAssignment") {
			interactions.UpdateAssignmentSubmit(bot, interactionCreate)
		} else {
			slog.Error(fmt.Sprintf("unknown modal submit: %s", interactionCreate.ModalSubmitData().CustomID))
		}
	default:
		slog.Error(fmt.Sprintf("unknown interaction type: %d", interactionCreate.Type))

	}
}
