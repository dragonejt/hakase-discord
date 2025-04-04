package interactions

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/views"
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

func SlashHakase(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate, hakaseClient clients.HakaseClient) {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(interactionCreate.ApplicationCommandData().Options))
	for _, opt := range interactionCreate.ApplicationCommandData().Options {
		optionMap[opt.Name] = opt
	}

	slog.Info(fmt.Sprintf("/hakase executed by %s (%s) in %s", interactionCreate.Member.User.Username, interactionCreate.Member.User.ID, interactionCreate.GuildID))
	transaction := sentry.StartTransaction(context.WithValue(context.Background(), clients.DiscordSession{}, bot), "/hakase")
	defer transaction.Finish()

	subcommand, exists := optionMap["cmd"]
	if !exists {
		ping(transaction, interactionCreate, hakaseClient)
	} else {
		switch subcommand.StringValue() {
		case "rock-paper-scissors":
			rockPaperScissors(transaction, interactionCreate)
		case "config":
			config(transaction, interactionCreate, hakaseClient)
		}
	}
}

func ping(span *sentry.Span, interactionCreate *discordgo.InteractionCreate, hakaseClient clients.HakaseClient) {
	span = span.StartChild("/hakase ping")
	defer span.Finish()
	bot := span.GetTransaction().Context().Value(clients.DiscordSession{}).(*discordgo.Session)

	start := time.Now()
	err := hakaseClient.HeadCourse(span, interactionCreate.GuildID)
	if err != nil {
		slog.Error(fmt.Sprintf("error pinging backend: %s", err.Error()))
	}

	err = bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("hakase pong! response time: %dms", time.Since(start).Milliseconds()),
		},
	})
	if err != nil {
		slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
	}

}

func rockPaperScissors(span *sentry.Span, interactionCreate *discordgo.InteractionCreate) {
	span = span.StartChild("/hakase rockPaperScissors")
	defer span.Finish()
	bot := span.GetTransaction().Context().Value(clients.DiscordSession{}).(*discordgo.Session)

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

func config(span *sentry.Span, interactionCreate *discordgo.InteractionCreate, hakaseClient clients.HakaseClient) {
	span = span.StartChild("/hakase config")
	defer span.Finish()
	bot := span.GetTransaction().Context().Value(clients.DiscordSession{}).(*discordgo.Session)

	course, err := hakaseClient.ReadCourse(span, interactionCreate.GuildID)
	if err != nil {
		slog.Error(fmt.Sprintf("error reading course: %s", err.Error()))
		err = bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("error reading course: %s", err.Error()),
			},
		})
		if err != nil {
			slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
		}
	}

	err = bot.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{views.ConfigView(course)},
			Components: views.ConfigActions(),
		},
	})
	if err != nil {
		slog.Error(fmt.Sprintf("error responding to interaction: %s", err.Error()))
	}
}
