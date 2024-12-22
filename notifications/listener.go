package notifications

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/dragonejt/hakase-discord/settings"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func ListenToStream(stopListener chan bool) {
	slog.Info(fmt.Sprintf("opening NATS consumer connection to: %s", settings.NATS_URL))
	connection, err := nats.Connect(settings.NATS_URL)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer func() {
		slog.Info("draining NATS consumer connection")
		err := connection.Drain()
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}()

	slog.Debug("opening jetstream consumer connection")
	js, err := jetstream.New(connection)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Debug(fmt.Sprintf("creating stream with name: %s", settings.STREAM_NAME))
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     settings.STREAM_NAME,
		Subjects: []string{fmt.Sprintf("%s.*", settings.STREAM_NAME)},
	})
	if err != nil {
		slog.Error(err.Error())
		return
	}

	consumer, err := js.CreateOrUpdateConsumer(ctx, settings.STREAM_NAME, jetstream.ConsumerConfig{
		Name:      settings.STREAM_NAME,
		Durable:   settings.STREAM_NAME,
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		slog.Error(err.Error())
		return
	}

	subscription, err := consumer.Consume(consumeMessage)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer subscription.Drain()

	<-stopListener

}

func consumeMessage(message jetstream.Msg) {
	slog.Info(fmt.Sprintf("received message: %s with subject: %s", string(message.Data()), message.Subject()))
	err := message.Ack()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}
