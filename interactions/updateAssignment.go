package interactions

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/views"
)

func UpdateAssignment(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	slog.Debug(fmt.Sprintf("updateAssignment executed by %s (%s) in %s", interactionCreate.Member.User.Username, interactionCreate.Member.User.ID, interactionCreate.GuildID))
	if interactionCreate.Member.Permissions&discordgo.PermissionAdministrator == 0 {
		err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "admin permissions needed!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
		}
		return
	}

	assignmentID := strings.Split(interactionCreate.MessageComponentData().CustomID, "_")[1]
	assignment, err := clients.ReadAssignment(assignmentID)
	if err != nil {
		slog.Error(fmt.Sprintf("error reading assignment: %s", err.Error()))
		err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: err.Error(),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
		}
	}
	err = bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   fmt.Sprintf("updateAssignment_%s", assignmentID),
			Title:      "update assignment",
			Components: views.AssignmentModal(&assignment),
		},
	})
	if err != nil {
		slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
	}
}

func UpdateAssignmentSubmit(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	slog.Info(fmt.Sprintf("updateAssignmentSubmit executed by %s (%s) in %s", interactionCreate.Member.User.Username, interactionCreate.Member.User.ID, interactionCreate.GuildID))
	err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
	}

	assignmentID := strings.Split(interactionCreate.ModalSubmitData().CustomID, "_")[1]
	assignmentData := interactionCreate.ModalSubmitData()
	assignment := clients.Assignment{
		CourseID: interactionCreate.GuildID,
	}

	if assignmentData.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value != "" {
		due, err := dateparse.ParseAny(assignmentData.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)
		if err != nil {
			_, err := bot.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
				Content: fmt.Sprintf("error parsing due date: %s", err.Error()),
			})
			if err != nil {
				slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
			}
			return
		}
		assignment.Due = due
	}

	if assignmentData.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value != "" {
		assignment.Name = assignmentData.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	}

	if assignmentData.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value != "" {
		assignment.Link = assignmentData.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	}

	currentAssignment, err := clients.ReadAssignment(assignmentID)
	if assignment.Due.Equal(time.Time{}) {
		assignment.Due = currentAssignment.Due
	} else if err == nil && assignment.Due.Before(currentAssignment.Due) {
		_, err := bot.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
			Content: "new due date before original assignment due date! hakase does not support this.",
		})
		if err != nil {
			slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
		}
		return
	}

	assignment.ID = currentAssignment.ID
	updatedAssignment, err := clients.UpdateAssignment(assignment)
	if err != nil {
		slog.Error(fmt.Sprintf("error updating assignment: %s", err.Error()))
		_, err := bot.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
			Content: err.Error(),
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
		}
		return
	}

	_, err = bot.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
		Content:    "assignment updated!",
		Embeds:     []*discordgo.MessageEmbed{views.AssignmentView(interactionCreate, updatedAssignment)},
		Components: []discordgo.MessageComponent{views.AssignmentActions(interactionCreate, updatedAssignment)},
	})
	if err != nil {
		slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
	}
}

func DeleteAssignment(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	slog.Debug(fmt.Sprintf("deleteAssignment executed by %s (%s) in %s", interactionCreate.Member.User.Username, interactionCreate.Member.User.ID, interactionCreate.GuildID))
	assignmentID := strings.Split(interactionCreate.MessageComponentData().CustomID, "_")[1]
	err := clients.DeleteAssignment(assignmentID)
	if err != nil {
		slog.Error(fmt.Sprintf("unable to delete assignment %s: %s", assignmentID, err.Error()))
		err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("unable to delete assignment %s: %s", assignmentID, err.Error()),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
		}
		return
	}

	err = bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("assignment %s deleted!", assignmentID),
		},
	})
	if err != nil {
		slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
	}
}
