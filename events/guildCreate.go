package events

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func GuildCreate(session *discordgo.Session, guildCreate *discordgo.GuildCreate) {
	slog.Info(fmt.Sprintf("joined guild: %s (%s)", guildCreate.Guild.Name, guildCreate.Guild.ID))
}
