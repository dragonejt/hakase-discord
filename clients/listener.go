// Package clients provides functions for consuming messages from NATS JetStream.
package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/palantir/stacktrace"
)

// ListenToStream starts a JetStream consumer and listens for messages, dispatching them to handlers.
func (mqClient *MQClient) ListenToStream(bot *discordgo.Session, hakaseClient BackendClient, stopListener chan bool) {
	slog.Info(fmt.Sprintf("opening NATS consumer connection to: %s", mqClient.NATSUrl))
	connection, err := nats.Connect(mqClient.NATSUrl)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error connecting to NATS: %s", mqClient.NATSUrl).Error())
		return
	}
	defer func() {
		slog.Info("draining NATS consumer connection")
		err := connection.Drain()
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error draining NATS connection").Error())
			return
		}
	}()

	slog.Debug("opening jetstream consumer connection")
	js, err := jetstream.New(connection)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error opening jetstream publisher connection").Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Debug(fmt.Sprintf("creating stream with name: %s", mqClient.StreamName))
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     mqClient.StreamName,
		Subjects: []string{fmt.Sprintf("%s.*", mqClient.StreamName)},
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error creating stream with name: %s", mqClient.StreamName).Error())
		return
	}

	consumer, err := js.CreateOrUpdateConsumer(ctx, mqClient.StreamName, jetstream.ConsumerConfig{
		Name:      mqClient.StreamName,
		Durable:   mqClient.StreamName,
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error creating consumer for stream: %s", mqClient.StreamName).Error())
		return
	}

	subscription, err := consumer.Consume(func(message jetstream.Msg) {
		consumeMessage(bot, hakaseClient, message)
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error subscribing to stream: %s", mqClient.StreamName).Error())
		return
	}
	defer subscription.Drain()

	<-stopListener

}

// consumeMessage dispatches messages based on their subject to the appropriate handler.
func consumeMessage(bot *discordgo.Session, hakaseClient BackendClient, message jetstream.Msg) {
	transaction := sentry.StartTransaction(context.WithValue(context.Background(), DiscordSession{}, bot), "consumeMessage")
	defer transaction.Finish()
	slog.Info(fmt.Sprintf("received message: %s with subject: %s", string(message.Data()), message.Subject()))

	if message.Subject() == "notifications" {
		consumeNotification(transaction, hakaseClient, message)
	} else if message.Subject() == "assignments" {
		consumeAssignmentNotification(transaction, hakaseClient, message)
	} else {
		slog.Error(fmt.Sprintf("unknown message subject: %s", message.Subject()))
		err := message.Ack()
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "failed to ACK message with subject: %s", message.Subject()).Error())
		}
	}

}

// consumeNotification handles notification messages received from JetStream.
func consumeNotification(span *sentry.Span, _ BackendClient, message jetstream.Msg) {
	span = span.StartChild("consumeNotification")
	defer span.Finish()

	slog.Info(fmt.Sprintf("received notification with message: %s", string(message.Data())))
}

// consumeAssignmentNotification handles assignment notification messages received from JetStream.
func consumeAssignmentNotification(span *sentry.Span, hakaseClient BackendClient, message jetstream.Msg) {
	span = span.StartChild("consumeAssignmentNotification")
	defer span.Finish()
	bot := span.GetTransaction().Context().Value(DiscordSession{}).(*discordgo.Session)

	assignmentNotification := AssignmentNotification{}
	err := json.Unmarshal(message.Data(), &assignmentNotification)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error unmarshalling assignment notification").Error())
		return
	}

	assignment, err := hakaseClient.ReadAssignment(span, fmt.Sprint(assignmentNotification.AssignmentID))
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to get assignment with ID: %d", assignmentNotification.AssignmentID).Error())
		return
	}

	course, err := hakaseClient.ReadCourse(span, assignmentNotification.CourseID)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to get course with courseID: %s", assignmentNotification.CourseID).Error())
		return
	}

	notificationsChannel := course.NotifyChannel
	if notificationsChannel == "" {
		guild, err := bot.Guild(course.CourseID)
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "unable to get guild system channel for notifications").Error())
			return
		}
		notificationsChannel = guild.SystemChannelID
	}

	notificationTime := assignment.Due.Add(-1 * assignmentNotification.Before)
	if time.Now().After(notificationTime) {
		_, err := bot.ChannelMessageSend(notificationsChannel, fmt.Sprintf("**[assignment notification]** assignment: %s is due in %s hours!", assignment.Name, assignmentNotification.Before/time.Hour))
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "failed to send assignment notification for %d", assignment.ID).Error())
			// retry sending assignment notification in 15 minutes
			_ = message.NakWithDelay(15 * time.Minute)
		}
	} else {
		err := message.NakWithDelay(time.Until(notificationTime))
		if err != nil {
			_, _ = bot.ChannelMessageSend(notificationsChannel, fmt.Sprintf("**[assignment notification error]** failed to schedule assignment notifications for assignment: %s", assignment.Name))
			slog.Error(stacktrace.Propagate(err, "failed to schedule assignment notifications for %d", assignment.ID).Error())
		}
	}
}
