package events

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func Ready(session *discordgo.Session, ready *discordgo.Ready) {
	slog.Info("logged in as " + ready.User.String())
}
