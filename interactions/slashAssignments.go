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

func SlashAssignments(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate, hakaseClient clients.HakaseClient) {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(interactionCreate.ApplicationCommandData().Options))
	for _, opt := range interactionCreate.ApplicationCommandData().Options {
		optionMap[opt.Name] = opt
	}

	slog.Info(fmt.Sprintf("/assignments executed by %s (%s) in %s", interactionCreate.Member.User.Username, interactionCreate.Member.User.ID, interactionCreate.GuildID))
	transaction := sentry.StartTransaction(context.WithValue(context.Background(), clients.DiscordSession{}, bot), "/assignments")
	defer transaction.Finish()

	assignmentID, exists := optionMap["id"]
	if exists {
		getAssignment(transaction, interactionCreate, hakaseClient, fmt.Sprint(assignmentID.IntValue()))

	} else {
		listAssignments(transaction, interactionCreate, hakaseClient)
	}

}

func getAssignment(span *sentry.Span, interactionCreate *discordgo.InteractionCreate, hakaseClient clients.HakaseClient, assignmentID string) {
	span = span.StartChild("/assignments getAssignment")
	defer span.Finish()
	bot := span.GetTransaction().Context().Value(clients.DiscordSession{}).(*discordgo.Session)
	assignment, err := hakaseClient.ReadAssignment(span, assignmentID)

	if err != nil {
		err = bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: err.Error(),
			},
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
	} else {
		err = bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds:     []*discordgo.MessageEmbed{views.AssignmentView(interactionCreate.Member, assignment)},
				Components: []discordgo.MessageComponent{views.AssignmentActions(assignment)},
			},
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
	}
}

func listAssignments(span *sentry.Span, interactionCreate *discordgo.InteractionCreate, hakaseClient clients.HakaseClient) {
	span = span.StartChild("/assignments listAssignments")
	defer span.Finish()
	bot := span.GetTransaction().Context().Value(clients.DiscordSession{}).(*discordgo.Session)
	assignments, err := hakaseClient.ListAssignments(span, interactionCreate.GuildID)

	if err != nil {
		err = bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: err.Error(),
			},
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
	} else {
		err = bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds:     []*discordgo.MessageEmbed{views.AssignmentsListView(interactionCreate.Member, assignments)},
				Components: []discordgo.MessageComponent{views.AssignmentsListActions()},
			},
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error responding to interaction").Error())
		}
	}
}
