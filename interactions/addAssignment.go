package interactions

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/araddon/dateparse"
	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/notifications"
	"github.com/dragonejt/hakase-discord/views"
)

func AddAssignment(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	slog.Debug(fmt.Sprintf("addAssignment executed by %s (%s)", interactionCreate.Member.User.Username, interactionCreate.Member.User.ID))
	err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   "addAssignment",
			Title:      "add assignment",
			Components: views.AssignmentModal(),
		},
	})
	if err != nil {
		slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
	}
}

func AddAssignmentSubmit(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	slog.Info(fmt.Sprintf("addAssignmentSubmit executed by %s (%s)", interactionCreate.Member.User.Username, interactionCreate.Member.User.ID))
	err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
		return
	}

	assignmentData := interactionCreate.ModalSubmitData()
	due, err := dateparse.ParseAny(assignmentData.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)
	if err != nil {
		slog.Error(fmt.Sprintf("error parsing due date: %s", err.Error()))
	}

	assignment := clients.Assignment{
		CourseID: interactionCreate.GuildID,
		Name:     assignmentData.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
		Due:      due,
		Link:     assignmentData.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
	}

	createdAssignment, err := clients.CreateAssignment(assignment)
	if err != nil {
		errMsg := err.Error()
		_, err = bot.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
			Content: errMsg,
		})
		if err != nil {
			slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
		}
		return
	}

	go notifications.PublishAssignmentNotification(notifications.AssignmentNotification{
		AssignmentID: createdAssignment.ID,
		Timestamp:    createdAssignment.Due.Add(-time.Hour * 24),
	})

	body, err := json.Marshal(createdAssignment)
	if err != nil {
		slog.Error(fmt.Sprintf("error marshalling assignment: %s", err.Error()))
		return
	}
	_, err = bot.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
		Content: string(body),
	})
	if err != nil {
		slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
	}
}
