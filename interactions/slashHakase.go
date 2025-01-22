package interactions

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/getsentry/sentry-go"
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

	slog.Info(fmt.Sprintf("/hakase executed by %s (%s) in %s", interactionCreate.Member.User.Username, interactionCreate.Member.User.ID, interactionCreate.GuildID))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	sentry.StartTransaction((ctx), "/hakase")

	subcommand, exists := optionMap["cmd"]
	if !exists {
		ping(bot, interactionCreate)
	} else {
		switch subcommand.StringValue() {
		case "rock-paper-scissors":
			rockPaperScissors(bot, interactionCreate)
		}
	}
}

func ping(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	start := time.Now()
	err := bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
	}
	_, err = clients.ReadCourse(interactionCreate.GuildID)
	if err != nil {
		slog.Error(fmt.Sprintf("error pinging backend: %s", err.Error()))
	}

	pong := fmt.Sprintf("hakase pong! response time: %dms", time.Since(start).Milliseconds())
	_, err = bot.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
		Content: &pong,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
	}

}

func rockPaperScissors(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
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
