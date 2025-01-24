package events

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
	"github.com/dragonejt/hakase-discord/clients"
)

type MockSession struct {
	mock.Mock
}

func (m *MockSession) UpdateCustomStatus(status string) error {
	args := m.Called(status)
	return args.Error(0)
}

func TestGuildCreate(t *testing.T) {
	mockSession := new(MockSession)
	mockGuildCreate := &discordgo.GuildCreate{
		Guild: &discordgo.Guild{
			ID:   "12345",
			Name: "Test Guild",
		},
	}

	mockSession.On("UpdateCustomStatus", "assisting 0 classes").Return(nil)

	GuildCreate(mockSession, mockGuildCreate)

	mockSession.AssertExpectations(t)
}

func TestGuildDelete(t *testing.T) {
	mockSession := new(MockSession)
	mockGuildDelete := &discordgo.GuildDelete{
		Guild: &discordgo.Guild{
			ID:   "12345",
			Name: "Test Guild",
		},
	}

	mockSession.On("UpdateCustomStatus", "assisting 0 classes").Return(nil)

	GuildDelete(mockSession, mockGuildDelete)

	mockSession.AssertExpectations(t)
}

func TestGuildCreateClientAPICall(t *testing.T) {
	mockSession := new(MockSession)
	mockGuildCreate := &discordgo.GuildCreate{
		Guild: &discordgo.Guild{
			ID:   "12345",
			Name: "Test Guild",
		},
	}

	mockSession.On("UpdateCustomStatus", "assisting 0 classes").Return(nil)

	GuildCreate(mockSession, mockGuildCreate)

	mockSession.AssertExpectations(t)

	// Assert that the correct clients API calls were called
	// This is a placeholder, replace with actual assertions based on your implementation
	// For example, if you have a mock client, you can assert that the expected methods were called
	// mockClient.AssertCalled(t, "ExpectedMethod", expectedArguments)
}

func TestGuildDeleteClientAPICall(t *testing.T) {
	mockSession := new(MockSession)
	mockGuildDelete := &discordgo.GuildDelete{
		Guild: &discordgo.Guild{
			ID:   "12345",
			Name: "Test Guild",
		},
	}

	mockSession.On("UpdateCustomStatus", "assisting 0 classes").Return(nil)

	GuildDelete(mockSession, mockGuildDelete)

	mockSession.AssertExpectations(t)

	// Assert that the correct clients API calls were called
	// This is a placeholder, replace with actual assertions based on your implementation
	// For example, if you have a mock client, you can assert that the expected methods were called
	// mockClient.AssertCalled(t, "ExpectedMethod", expectedArguments)
}
