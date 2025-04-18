package views

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
)

func AssignmentView(member *discordgo.Member, assignment clients.Assignment) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       assignment.Name,
		Description: fmt.Sprintf("due %s", assignment.Due.Format(time.RFC1123)),
		Author:      &discordgo.MessageEmbedAuthor{Name: member.User.Username, IconURL: member.User.AvatarURL("")},
		URL:         assignment.Link,
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("id %d", assignment.ID)},
	}
}

func AssignmentActions(assignment clients.Assignment) *discordgo.ActionsRow {
	return &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Emoji: &discordgo.ComponentEmoji{
					Name: "📝",
				},
				Label:    "edit",
				Style:    discordgo.PrimaryButton,
				CustomID: fmt.Sprintf("updateAssignmentAction_%d", assignment.ID),
			},
			discordgo.Button{
				Emoji: &discordgo.ComponentEmoji{
					Name: "🗑️",
				},
				Label:    "remove",
				Style:    discordgo.SecondaryButton,
				CustomID: fmt.Sprintf("deleteAssignmentAction_%d", assignment.ID),
			},
		},
	}
}

func AssignmentModal(assignment *clients.Assignment) []discordgo.MessageComponent {

	newAssignment := assignment == nil

	if newAssignment {
		assignment = &clients.Assignment{
			// placeholder data
			Name: "Assignment 1",
			Due:  time.Now(),
			Link: "https://canvas.instructure.com",
		}
	}
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "assignmentName",
					Label:       "assignment name:",
					Style:       discordgo.TextInputShort,
					Placeholder: assignment.Name,
					Required:    newAssignment,
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
					Placeholder: assignment.Due.Format(time.RFC1123),
					Required:    newAssignment,
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
					Placeholder: assignment.Link,
					Required:    false,
					MaxLength:   50,
				},
			},
		},
	}
}
