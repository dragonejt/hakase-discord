package interactions

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/views"
	"github.com/getsentry/sentry-go"
)

var AssignmentsCommand = discordgo.ApplicationCommand{
	Name:        "assignments",
	Description: "configure assignments for due date notifications",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "id",
			Description: "retrieves assignment with this id",
			Type:        discordgo.ApplicationCommandOptionInteger,
		},
	},
}

func SlashAssignments(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
		return
	}

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(interactionCreate.ApplicationCommandData().Options))
	for _, opt := range interactionCreate.ApplicationCommandData().Options {
		optionMap[opt.Name] = opt
	}

	slog.Info(fmt.Sprintf("/assignments executed by %s (%s) in %s", interactionCreate.Member.User.Username, interactionCreate.Member.User.ID, interactionCreate.GuildID))
	sentry.StartTransaction(context.Background(), "/assignments")

	assignmentID, exists := optionMap["id"]
	if exists {
		getAssignment(bot, interactionCreate, fmt.Sprint(assignmentID.IntValue()))

	} else {
		listAssignments(bot, interactionCreate)
	}

}

func getAssignment(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate, assignmentID string) {
	assignment, err := clients.ReadAssignment(assignmentID)

	if err != nil {
		errMsg := err.Error()
		_, err = bot.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		if err != nil {
			slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
		}
	} else {
		_, err = bot.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
			Embeds:     &[]*discordgo.MessageEmbed{views.AssignmentView(interactionCreate, assignment)},
			Components: &[]discordgo.MessageComponent{views.AssignmentActions(interactionCreate, assignment)},
		})
		if err != nil {
			slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
		}
	}
}

func listAssignments(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	assignments, err := clients.ListAssignments(interactionCreate.GuildID)

	if err != nil {
		errMsg := err.Error()
		_, err = bot.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		if err != nil {
			slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
		}
	} else {
		_, err = bot.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
			Embeds:     &[]*discordgo.MessageEmbed{views.AssignmentsListView(interactionCreate, assignments)},
			Components: &[]discordgo.MessageComponent{views.AssignmentsListActions(interactionCreate)},
		})
		if err != nil {
			slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
		}
	}
}
