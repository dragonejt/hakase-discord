package commands

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
)

var AssignmentsCommand *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
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
	bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(interactionCreate.ApplicationCommandData().Options))
	for _, opt := range interactionCreate.ApplicationCommandData().Options {
		optionMap[opt.Name] = opt
	}

	slog.Info(fmt.Sprintf("/assignments %s", fmt.Sprint(optionMap)))

	assignmentID, exists := optionMap["id"]
	if exists {
		slog.Info(assignmentID.StringValue())

	} else {
		assignments, err := clients.ListAssignments(interactionCreate.GuildID)
		if err != nil {
			bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: err.Error(),
				},
			})
		} else {
			body, err := json.Marshal(assignments)
			if err != nil {
				slog.Error(err.Error())
				return
			}
			bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: string(body),
				},
			})
		}
	}

}
