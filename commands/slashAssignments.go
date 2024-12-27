package commands

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
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

	slog.Info(fmt.Sprintf("/assignments %s", fmt.Sprint(optionMap)))

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
		err = bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: err.Error(),
			},
		})
		if err != nil {
			slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
		}
	} else {
		body, err := json.Marshal(assignments)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		err = bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: string(body),
			},
		})
		if err != nil {
			slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
		}
	}
}
