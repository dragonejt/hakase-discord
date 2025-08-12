package main

import (
	"fmt"
	"log"
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
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if !settings.DEBUG {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              "https://1ef54abfefe4da3ed56195664ee3fc03@o4507124907638784.ingest.us.sentry.io/4509828350476288",
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

	} else {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	bot, err := discordgo.New(fmt.Sprintf("Bot %s", settings.DISCORD_BOT_TOKEN))
	if err != nil {
		slog.Error(fmt.Sprintf("error creating discord session: %s", err.Error()))
		return
	}
	bot.StateEnabled = true

	backend := clients.BackendClient{
		URL:         settings.BACKEND_URL,
		API_KEY:     settings.BACKEND_API_KEY,
		HTTP_CLIENT: bot.Client,
	}

	stopListener := make(chan bool, 1)
	go clients.ListenToStream(bot, stopListener)

	err = bot.Open()
	if err != nil {
		slog.Error(fmt.Sprintf("error opening connection to discord: %s", err.Error()))
		return
	}
	defer bot.Close()

	slog.Info("registering event handlers")
	bot.AddHandler(events.Ready)
	bot.AddHandler(func(bot *discordgo.Session, guildCreate *discordgo.GuildCreate) {
		events.GuildCreate(bot, guildCreate, &backend)
	})
	bot.AddHandler(func(bot *discordgo.Session, guildDelete *discordgo.GuildDelete) {
		events.GuildDelete(bot, guildDelete, &backend)
	})
	bot.AddHandler(func(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
		events.InteractionCreate(bot, interactionCreate, &backend)
	})

	slog.Info("registering interactions")
	interactions := []*discordgo.ApplicationCommand{&interactions.AssignmentsCommand, &interactions.HakaseCommand}
	for _, cmd := range interactions {
		_, err = bot.ApplicationCommandCreate(bot.State.User.ID, "", cmd)
		if err != nil {
			slog.Error(fmt.Sprintf("error registering command: %s, %s", cmd.Name, err.Error()))
		} else {
			slog.Info(fmt.Sprintf("successfully registered command: %s", cmd.Name))
		}
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	stopListener <- true
}
