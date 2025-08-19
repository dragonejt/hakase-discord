// Package interactions provides handlers for assignment list actions (add assignment).
package interactions

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/araddon/dateparse"
	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/views"
	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
)

// AddAssignment opens a modal for adding a new assignment via Discord interaction.
func AddAssignment(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	transaction := sentry.StartTransaction(context.WithValue(context.Background(), clients.DiscordSession{}, bot), "addAssignmentAction")
	defer transaction.Finish()
	slog.Debug(fmt.Sprintf("addAssignment executed by %s (%s) in %s", interactionCreate.Member.User.Username, interactionCreate.Member.User.ID, interactionCreate.GuildID))
	if interactionCreate.Member.Permissions&discordgo.PermissionAdministrator == 0 {
		err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "admin permissions needed!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
		return
	}

	err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   "addAssignment",
			Title:      "add assignment",
			Components: views.AssignmentModal(nil),
		},
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
	}

}

// AddAssignmentSubmit handles the submission of the add assignment modal and creates the assignment.
func AddAssignmentSubmit(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate, hakaseClient clients.HakaseClient) {
	slog.Info(fmt.Sprintf("addAssignmentSubmit executed by %s (%s) in %s", interactionCreate.Member.User.Username, interactionCreate.Member.User.ID, interactionCreate.GuildID))
	transaction := sentry.StartTransaction(context.WithValue(context.Background(), clients.DiscordSession{}, bot), "addAssignmentSubmit")
	defer transaction.Finish()

	err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
	}

	assignmentData := interactionCreate.ModalSubmitData()
	due, err := dateparse.ParseAny(assignmentData.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)
	if err != nil {
		_, err := bot.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
			Content: fmt.Sprintf("error parsing due date: %s", err.Error()),
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
		return
	}

	assignment := clients.Assignment{
		Name:     assignmentData.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
		Due:      due,
		CourseID: interactionCreate.GuildID,
	}

	if assignmentData.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value != "" {
		assignment.Link = assignmentData.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	}

	if assignment.Due.Before(time.Now()) {
		_, err := bot.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
			Content: "due date before current time! hakase does not support this.",
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
		return
	}

	createdAssignment, err := hakaseClient.CreateAssignment(transaction, assignment)
	if err != nil {
		_, err := bot.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
			Content: err.Error(),
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
		return
	}

	go clients.PublishAssignmentNotification(transaction, clients.AssignmentNotification{
		AssignmentID: createdAssignment.ID,
		CourseID:     interactionCreate.GuildID,
		Before:       time.Hour,
	})

	go clients.PublishAssignmentNotification(transaction, clients.AssignmentNotification{
		AssignmentID: createdAssignment.ID,
		CourseID:     interactionCreate.GuildID,
		Before:       time.Hour * 24,
	})

	_, err = bot.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
		Content:    "assignment created!",
		Embeds:     []*discordgo.MessageEmbed{views.AssignmentView(interactionCreate.Member, createdAssignment)},
		Components: []discordgo.MessageComponent{views.AssignmentActions(createdAssignment)},
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
	}
}
