// hakase-discord is the entry point for the Discord bot.
// It initializes logging, Sentry, Discord session, backend client, and event handlers.
// It registers application commands and starts the bot event loop.
package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/events"
	"github.com/dragonejt/hakase-discord/interactions"
	"github.com/dragonejt/hakase-discord/settings"
	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
)

func main() {
	if settings.DEBUG {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	if settings.SENTRY_DSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              settings.SENTRY_DSN,
			EnableTracing:    true,
			SampleRate:       1,
			TracesSampleRate: 1,
			SendDefaultPII:   true,
			EnableLogs:       true,
			Environment:      settings.ENV,
		})
		if err != nil {
			slog.Warn(fmt.Sprintf("error initiating sentry: %s", err))
		}
		slog.SetDefault(slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stderr, sentry.NewLogger(context.Background())), &slog.HandlerOptions{AddSource: true})))
	}

	bot, err := discordgo.New(fmt.Sprintf("Bot %s", settings.DISCORD_BOT_TOKEN))
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to create discord session").Error())
		return
	}
	bot.StateEnabled = true

	backend := &clients.BackendClient{
		URL:         settings.BACKEND_URL,
		API_KEY:     settings.BACKEND_API_KEY,
		HTTP_CLIENT: bot.Client,
	}

	stopListener := make(chan bool, 1)
	go clients.ListenToStream(bot, backend, stopListener)

	err = bot.Open()
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to open discord session").Error())
		return
	}

	slog.Info("registering event handlers")
	bot.AddHandler(events.Ready)
	bot.AddHandler(func(bot *discordgo.Session, guildCreate *discordgo.GuildCreate) {
		events.GuildCreate(bot, guildCreate, backend)
	})
	bot.AddHandler(func(bot *discordgo.Session, guildDelete *discordgo.GuildDelete) {
		events.GuildDelete(bot, guildDelete, backend)
	})
	bot.AddHandler(func(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
		events.InteractionCreate(bot, interactionCreate, backend)
	})

	slog.Info("registering interactions")
	interactions := []*discordgo.ApplicationCommand{&interactions.AssignmentsCommand, &interactions.HakaseCommand}
	for _, cmd := range interactions {
		_, err = bot.ApplicationCommandCreate(bot.State.User.ID, "", cmd)
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "failed to register command: %s", cmd.Name).Error())
		} else {
			slog.Info(fmt.Sprintf("successfully registered command: %s", cmd.Name))
		}
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	stopListener <- true

	err = bot.Close()
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to close discord session").Error())
	}
}
