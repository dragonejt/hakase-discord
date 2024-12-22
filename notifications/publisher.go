package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/dragonejt/hakase-discord/settings"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var PublisherPool sync.Pool = sync.Pool{
	New: createStreamConnection,
}

func createStreamConnection() any {
	slog.Info(fmt.Sprintf("opening NATS publisher connection to: %s", settings.NATS_URL))
	connection, err := nats.Connect(settings.NATS_URL)
	if err != nil {
		panic(err.Error())
	}

	slog.Debug("opening jetstream publisher connection")
	js, err := jetstream.New(connection)
	if err != nil {
		panic(err.Error())
	}
	return js
}

func publishMessage(subject string, message []byte) error {
	js := PublisherPool.Get().(jetstream.JetStream)
	defer PublisherPool.Put(js)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Debug(fmt.Sprintf("publishing message to subject: %s.%s", settings.STREAM_NAME, subject))
	_, err := js.Publish(ctx, fmt.Sprintf("%s.%s", settings.STREAM_NAME, subject), message)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

func PublishNotification(notification string) {
	err := publishMessage("notifications", []byte(notification))
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}

func PublishAssignmentNotification(notification AssignmentNotification) {
	message, err := json.Marshal(notification)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
	err = publishMessage("assignments", message)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}

func PublishStudySessionNotification(notification StudySessionNotification) {
	message, err := json.Marshal(notification)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	err = publishMessage("study_sessions", message)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}
