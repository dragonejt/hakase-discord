// Package views provides Discord message embeds and components for course configuration.
package views

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
)

// ConfigView returns a Discord message embed displaying the configuration for a course.
// It shows the notifications channel and role for the given course.
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

// ConfigActions returns Discord message components for updating the course's notifications channel and role.
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
