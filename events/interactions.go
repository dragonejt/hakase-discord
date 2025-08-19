// Package events provides the Discord interaction event handler.
package events

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/interactions"
)

// InteractionCreate dispatches Discord interactions to the appropriate handler based on type and command.
func InteractionCreate(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate, hakaseClient clients.HakaseClient) {
	switch interactionCreate.Type {
	case discordgo.InteractionApplicationCommand:
		switch interactionCreate.ApplicationCommandData().Name {
		case "assignments":
			interactions.SlashAssignments(bot, interactionCreate, hakaseClient)
		case "hakase":
			interactions.SlashHakase(bot, interactionCreate, hakaseClient)
		default:
			slog.Error(fmt.Sprintf("unknown command: %s", interactionCreate.ApplicationCommandData().Name))
		}
	case discordgo.InteractionMessageComponent:
		customID := interactionCreate.MessageComponentData().CustomID
		if strings.HasPrefix(customID, "addAssignmentAction") {
			interactions.AddAssignment(bot, interactionCreate)
		} else if strings.HasPrefix(customID, "updateAssignmentAction") {
			interactions.UpdateAssignment(bot, interactionCreate, hakaseClient)
		} else if strings.HasPrefix(customID, "deleteAssignmentAction") {
			interactions.DeleteAssignment(bot, interactionCreate, hakaseClient)
		} else if strings.HasPrefix(customID, "updateNotifyChannel") {
			interactions.UpdateNotifyChannel(bot, interactionCreate, hakaseClient)
		} else if strings.HasPrefix(customID, "updateNotifyRole") {
			interactions.UpdateNotifyRole(bot, interactionCreate, hakaseClient)
		} else {
			slog.Error(fmt.Sprintf("unknown message component action: %s", customID))
		}
	case discordgo.InteractionModalSubmit:
		customID := interactionCreate.ModalSubmitData().CustomID
		if strings.HasPrefix(customID, "addAssignment") {
			interactions.AddAssignmentSubmit(bot, interactionCreate, hakaseClient)
		} else if strings.HasPrefix(customID, "updateAssignment") {
			interactions.UpdateAssignmentSubmit(bot, interactionCreate, hakaseClient)
		} else {
			slog.Error(fmt.Sprintf("unknown modal submit: %s", interactionCreate.ModalSubmitData().CustomID))
		}
	default:
		slog.Error(fmt.Sprintf("unknown interaction type: %d", interactionCreate.Type))

	}
}
