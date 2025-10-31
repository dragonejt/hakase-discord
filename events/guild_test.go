package events_test

import (
	"log/slog"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/events"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GuildEventsTestSuite struct {
	suite.Suite
	bot           *discordgo.Session
	guildCreate   *discordgo.GuildCreate
	guildDelete   *discordgo.GuildDelete
	hakaseClient  *clients.HakaseClient
	backendClient *MockBackendClient
}

func TestGuildEvents(t *testing.T) {
	suite.Run(t, new(GuildEventsTestSuite))
}

type MockBackendClient struct {
	clients.BackendClient
	mock.Mock
}

func (hakaseClient *MockBackendClient) CreateCourse(span *sentry.Span, course clients.Course) error {
	hakaseClient.Called(span, course)
	return nil
}

func (hakaseClient *MockBackendClient) DeleteCourse(span *sentry.Span, courseID string) error {
	hakaseClient.Called(span, courseID)
	return nil
}

func (testSuite *GuildEventsTestSuite) SetupTest() {
	testSuite.bot = &discordgo.Session{}
	testSuite.bot.State = &discordgo.State{}
	guild := &discordgo.Guild{
		ID:   "1234567890",
		Name: "Test Guild",
	}
	testSuite.bot.State.Guilds = []*discordgo.Guild{guild}
	testSuite.guildCreate = &discordgo.GuildCreate{
		Guild: guild,
	}
	testSuite.guildDelete = &discordgo.GuildDelete{
		Guild:        guild,
		BeforeDelete: guild,
	}
	testSuite.backendClient = new(MockBackendClient)
	testSuite.hakaseClient = &clients.HakaseClient{
		Backend:       testSuite.backendClient,
		Notifications: nil,
	}
}

func (testSuite *GuildEventsTestSuite) TestGuildCreateSuccess() {
	testSuite.backendClient.On("CreateCourse", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		slog.Info("CreateCourse called")
	})
	events.GuildCreate(testSuite.bot, testSuite.guildCreate, *testSuite.hakaseClient)
	testSuite.backendClient.AssertExpectations(testSuite.T())
}

func (testSuite *GuildEventsTestSuite) TestGuildDeleteSuccess() {
	testSuite.backendClient.On("DeleteCourse", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		slog.Info("DeleteCourse called")
	})
	events.GuildDelete(testSuite.bot, testSuite.guildDelete, *testSuite.hakaseClient)
	testSuite.backendClient.AssertExpectations(testSuite.T())
}
