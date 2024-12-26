package commands

import (
	"fmt"
	"log/slog"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

var HakaseCommand = discordgo.ApplicationCommand{
	Name:        "hakase",
	Description: "course configuration",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "cmd",
			Description: "subcommand to execute",
			Type:        discordgo.ApplicationCommandOptionString,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "rock-paper-scissors",
					Value: "rock-paper-scissors",
				},
				{
					Name:  "config",
					Value: "config",
				},
			},
		},
	},
}

var rockPaperScissorsGIFS = []string{
	"https://tenor.com/view/hakase-rock-paper-scissors-nichijou-gif-16268712187530688616",
	"https://tenor.com/view/nichijou-hakase-rps-rock-paper-scissors-nano-gif-17854309283562565671",
	"https://tenor.com/view/rps-nichijou-hakase-nano-rock-paper-scissors-gif-9851171842395079248",
	"https://tenor.com/view/hakase-rock-paper-scissors-nichijou-gif-8852988661228140380",
	"https://tenor.com/view/rps-nichijou-hakase-nano-rock-paper-scissors-gif-9851171842395079248",
	"https://tenor.com/view/nichijou-hakase-rps-rock-paper-scissors-nano-gif-11850067363499322337",
}

func SlashHakase(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(interactionCreate.ApplicationCommandData().Options))
	for _, opt := range interactionCreate.ApplicationCommandData().Options {
		optionMap[opt.Name] = opt
	}

	subcommand, exists := optionMap["subcommand"]
	if !exists {
		err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "hakase pong!",
			},
		})
		if err != nil {
			slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
		}
	} else {
		if subcommand.StringValue() == "rock-paper-scissors" {
			err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: rockPaperScissorsGIFS[rand.Intn(len(rockPaperScissorsGIFS))],
				},
			})
			if err != nil {
				slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
			}
		}
	}
}
