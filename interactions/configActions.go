// Package interactions provides handlers for course config actions (update notify channel/role).
package interactions

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/views"
	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
)

// UpdateNotifyChannel updates the notifications channel for a course based on user interaction.
func UpdateNotifyChannel(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate, hakaseClient clients.HakaseClient) {
	transaction := sentry.StartTransaction(context.WithValue(context.Background(), clients.DiscordSession{}, bot), "updateNotifyChannel")
	defer transaction.Finish()
	slog.Debug(fmt.Sprintf("updateNotifyChannel executed by %s (%s) in %s", interactionCreate.Member.User.Username, interactionCreate.Member.User.ID, interactionCreate.GuildID))
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

	notifyChannel := interactionCreate.MessageComponentData().Values[0]
	err := hakaseClient.UpdateCourse(transaction, clients.Course{
		CourseID:      interactionCreate.GuildID,
		NotifyChannel: notifyChannel,
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error updating course").Error())
		err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("error updating course: %s", err.Error()),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
	}

	updatedCourse, err := hakaseClient.ReadCourse(transaction, interactionCreate.GuildID)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error reading updated course").Error())
		err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("error reading updated course course: %s", err.Error()),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
	}

	err = bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "notifications channel updated!",
			Embeds:  []*discordgo.MessageEmbed{views.ConfigView(updatedCourse)},
		},
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
	}
}

// UpdateNotifyRole updates the notifications role for a course based on user interaction.
func UpdateNotifyRole(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate, hakaseClient clients.HakaseClient) {
	transaction := sentry.StartTransaction(context.WithValue(context.Background(), clients.DiscordSession{}, bot), "updateNotifyRole")
	defer transaction.Finish()
	slog.Debug(fmt.Sprintf("updateNotifyRole executed by %s (%s) in %s", interactionCreate.Member.User.Username, interactionCreate.Member.User.ID, interactionCreate.GuildID))
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

	notifyRole := interactionCreate.MessageComponentData().Values[0]
	err := hakaseClient.UpdateCourse(transaction, clients.Course{
		CourseID:    interactionCreate.GuildID,
		NotifyGroup: notifyRole,
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error updating course").Error())
		err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("error updating course: %s", err.Error()),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
	}

	updatedCourse, err := hakaseClient.ReadCourse(transaction, interactionCreate.GuildID)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error reading updated course").Error())
		err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("error reading updated course course: %s", err.Error()),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
	}

	err = bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "notifications role updated!",
			Embeds:  []*discordgo.MessageEmbed{views.ConfigView(updatedCourse)},
		},
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
	}
}
