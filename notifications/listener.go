package notifications

import (
	"log/slog"
	"os"
	"strings"

	"github.com/nats-io/nats.go"
)

func ListenToStream() {
	natsURL := os.Getenv("NATS_URL")
	streamName := os.Getenv("STREAM_NAME")
	slog.Debug("opening NATS consumer connection to: " + natsURL)
	connection, err := nats.Connect(natsURL)
	if err != nil {
		panic(err.Error())
	}
	defer connection.Drain()

	slog.Debug("opening jetstream consumer connection")
	jsctx, err := connection.JetStream()
	if err != nil {
		panic(err)
	}

	slog.Debug("creating NATS stream")
	stream, err := jsctx.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{streamName + ".*"},
	})
	if err != nil {
		if !strings.Contains(err.Error(), "stream name already in use") {

			slog.Info("found stream with name: " + stream.Config.Name + ", updating it")
			_, err = jsctx.UpdateStream(&nats.StreamConfig{
				Name:     streamName,
				Subjects: []string{streamName + ".*"},
			})
			if err != nil {
				slog.Error(err.Error())
				panic(err)
			}
		} else {
			slog.Error(err.Error())
			panic(err)
		}
	}

	subscription, err := jsctx.Subscribe(streamName+".notifications", consumeMessage, nats.BindStream(stream.Config.Name))
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	slog.Info("subscription created to: " + subscription.Subject)

	select {}

}

func consumeMessage(message *nats.Msg) {
	slog.Info("Received Message: " + string(message.Data))
	err := message.Ack()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}
