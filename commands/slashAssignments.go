package commands

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/views"
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

	slog.Info(fmt.Sprintf("/assignments executed by %s (%s)", interactionCreate.Member.User.Username, interactionCreate.Member.User.ID))

	assignmentID, exists := optionMap["id"]
	if exists {
		slog.Info(assignmentID.StringValue())

	} else {
		listAssignments(bot, interactionCreate)
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
			Components: &[]discordgo.MessageComponent{views.AssignmentsListActions(interactionCreate, assignments)},
		})
		if err != nil {
			slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
		}
	}
}
