package events

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
	"github.com/dragonejt/hakase-discord/interactions"
)

type MockSession struct {
	mock.Mock
}

func (m *MockSession) UpdateCustomStatus(status string) error {
	args := m.Called(status)
	return args.Error(0)
}

func TestInteractionCreate(t *testing.T) {
	mockSession := new(MockSession)
	mockInteractionCreate := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID: "12345",
		},
	}

	mockSession.On("UpdateCustomStatus", "assisting 0 classes").Return(nil)

	InteractionCreate(mockSession, mockInteractionCreate)

	mockSession.AssertExpectations(t)
}
