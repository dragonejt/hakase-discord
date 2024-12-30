package views

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
)

func AssignmentsListView(interactionCreate *discordgo.InteractionCreate, assignments []clients.Assignment) *discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{
		Title:       "assignments",
		Description: fmt.Sprintf("%d assignments in course", len(assignments)),
		Author:      &discordgo.MessageEmbedAuthor{Name: interactionCreate.Member.User.Username, IconURL: interactionCreate.Member.User.AvatarURL("")},
	}

	for _, assignment := range assignments {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("id: %d", assignment.ID),
			Value:  assignment.Name,
			Inline: true,
		})
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "due:",
			Value:  assignment.Due.Format(time.RFC1123),
			Inline: true,
		})
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "\u200B",
			Value:  "\u200B",
			Inline: false,
		})
	}

	return &embed

}

func AssignmentsListActions(interactionCreate *discordgo.InteractionCreate) *discordgo.ActionsRow {
	return &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Emoji: &discordgo.ComponentEmoji{
					Name: "➕",
				},
				Label:    "create",
				Style:    discordgo.PrimaryButton,
				CustomID: "addAssignmentAction",
			},
		},
	}

}

func AssignmentView(interactionCreate *discordgo.InteractionCreate, assignment clients.Assignment) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       assignment.Name,
		Description: fmt.Sprintf("due %s", assignment.Due.Format(time.RFC1123)),
		Author:      &discordgo.MessageEmbedAuthor{Name: interactionCreate.Member.User.Username, IconURL: interactionCreate.Member.User.AvatarURL("")},
		URL:         assignment.Link,
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("id %d", assignment.ID)},
	}
}

func AssignmentActions(interactionCreate *discordgo.InteractionCreate) *discordgo.ActionsRow {
	return &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Emoji: &discordgo.ComponentEmoji{
					Name: "📝",
				},
				Label:    "update",
				Style:    discordgo.PrimaryButton,
				CustomID: "updateAssignmentAction",
			},
			discordgo.Button{
				Emoji: &discordgo.ComponentEmoji{
					Name: "🗑️",
				},
				Label:    "delete",
				Style:    discordgo.SecondaryButton,
				CustomID: "deleteAssignmentAction",
			},
		},
	}
}

func AssignmentModal() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "assignmentName",
					Label:       "assignment name:",
					Style:       discordgo.TextInputShort,
					Placeholder: "Assignment 1",
					Required:    true,
					MaxLength:   50,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "assignmentDue",
					Label:       "due date:",
					Style:       discordgo.TextInputShort,
					Placeholder: time.Now().Format(time.RFC1123),
					Required:    true,
					MaxLength:   50,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "assignmentLink",
					Label:       "link:",
					Style:       discordgo.TextInputShort,
					Placeholder: "https://canvas.instructure.com",
					Required:    false,
					MaxLength:   50,
				},
			},
		},
	}
}
