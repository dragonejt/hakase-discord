// Package clients provides functions for publishing notifications to NATS JetStream.
package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/palantir/stacktrace"
)

func CreateStreamConnection(NATS_URL string) jetstream.JetStream {
	slog.Info(fmt.Sprintf("opening NATS publisher connection to: %s", NATS_URL))
	connection, err := nats.Connect(NATS_URL)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error connecting to NATS: %s", NATS_URL).Error())
		return nil
	}

	slog.Debug("opening jetstream publisher connection")
	js, err := jetstream.New(connection)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error opening jetstream publisher connection").Error())
		return nil
	}
	return js
}

func (mqClient *MQClient) publishMessage(span *sentry.Span, subject string, message []byte) error {
	js := mqClient.PublisherPool.Get().(jetstream.JetStream)
	defer mqClient.PublisherPool.Put(js)

	ctx, cancel := context.WithTimeout(span.Context(), 10*time.Second)
	defer cancel()

	slog.Debug(fmt.Sprintf("publishing message to subject: %s.%s", mqClient.StreamName, subject))
	_, err := js.Publish(ctx, fmt.Sprintf("%s.%s", mqClient.StreamName, subject), message)
	if err != nil {
		return stacktrace.Propagate(err, "error publishing message to subject: %s.%s", mqClient.StreamName, subject)
	}

	return nil
}

// PublishNotification publishes a notification message to the notifications subject in JetStream.
func (mqClient *MQClient) PublishNotification(span *sentry.Span, notification string) {
	span = span.StartChild("publishNotification")
	defer span.Finish()

	err := mqClient.publishMessage(span, "notifications", []byte(notification))
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error publishing notification").Error())
		return
	}
}

// PublishAssignmentNotification publishes an assignment notification to the assignments subject in JetStream.
func (mqClient *MQClient) PublishAssignmentNotification(span *sentry.Span, notification AssignmentNotification) {
	span = span.StartChild("publishAssignmentNotification")
	defer span.Finish()

	message, err := json.Marshal(notification)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error marshalling assignment notification").Error())
		return
	}
	err = mqClient.publishMessage(span, "assignments", message)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error publishing assignment notification").Error())
		return
	}
}

// PublishStudySessionNotification publishes a study session notification to the study_sessions subject in JetStream.
func (mqClient *MQClient) PublishStudySessionNotification(span *sentry.Span, notification StudySessionNotification) {
	span = span.StartChild("publishStudySessionNotification")
	defer span.Finish()

	message, err := json.Marshal(notification)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error marshalling study session notification").Error())
		return
	}

	err = mqClient.publishMessage(span, "study_sessions", message)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error publishing study session notification").Error())
		return
	}
}
