package views

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
)

func ConfigView(course clients.Course) *discordgo.MessageEmbed {
	notifyChannel, notifyRole := course.NotifyChannel, course.NotifyGroup
	if notifyChannel != "" {
		notifyChannel = fmt.Sprintf("<#%s>", notifyChannel)
	}
	if notifyRole != "" {
		notifyRole = fmt.Sprintf("<@&%s>", notifyRole)
	}

	return &discordgo.MessageEmbed{
		Title: "course config",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "notifications channel",
				Value: notifyChannel,
			},
			{
				Name:  "notifications role",
				Value: notifyRole,
			},
		},
	}
}

func ConfigActions() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					MenuType:    discordgo.ChannelSelectMenu,
					CustomID:    "updateNotifyChannel",
					Placeholder: "update notifications channel",
				},
			},
		},
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					MenuType:    discordgo.RoleSelectMenu,
					CustomID:    "updateNotifyRole",
					Placeholder: "update notifications role",
				},
			},
		},
	}
}
