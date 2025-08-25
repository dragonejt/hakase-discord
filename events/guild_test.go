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
	bot          *discordgo.Session
	guildCreate  *discordgo.GuildCreate
	guildDelete  *discordgo.GuildDelete
	hakaseClient *MockHakaseClient
}

func TestGuildEvents(t *testing.T) {
	suite.Run(t, new(GuildEventsTestSuite))
}

type MockHakaseClient struct {
	clients.HakaseClient
	mock.Mock
}

func (hakaseClient *MockHakaseClient) CreateCourse(span *sentry.Span, course clients.Course) error {
	hakaseClient.Called(span, course)
	return nil
}

func (hakaseClient *MockHakaseClient) DeleteCourse(span *sentry.Span, courseID string) error {
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
	testSuite.hakaseClient = new(MockHakaseClient)
}

func (testSuite *GuildEventsTestSuite) TestGuildCreateSuccess() {
	testSuite.hakaseClient.On("CreateCourse", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		slog.Info("CreateCourse called")
	})
	events.GuildCreate(testSuite.bot, testSuite.guildCreate, testSuite.hakaseClient)
	testSuite.hakaseClient.AssertExpectations(testSuite.T())
}

func (testSuite *GuildEventsTestSuite) TestGuildDeleteSuccess() {
	testSuite.hakaseClient.On("DeleteCourse", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		slog.Info("DeleteCourse called")
	})
	events.GuildDelete(testSuite.bot, testSuite.guildDelete, testSuite.hakaseClient)
	testSuite.hakaseClient.AssertExpectations(testSuite.T())
}
