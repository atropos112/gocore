package events

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/atropos112/gocore/utils"
	nc "github.com/nats-io/nats.go"
)

// NatsUrlPrefix  is the prefix for NATS URLs.
const NatsURLPrefix = "nats://"

// CreateSubscriber creates a new NATS publisher.
// Can then use it to subscribe to topic like
// messages, err := subscriber.Subscribe(context.Background(), "example.topic")
func CreateNATSSubscriber(natsUrl string) (*nats.Subscriber, error) {
	if natsUrl[:len(NatsURLPrefix)] != NatsURLPrefix {
		panic("natsUrl must start with " + NatsURLPrefix)
	}

	options := []nc.Option{
		nc.RetryOnFailedConnect(true),
		nc.Timeout(30 * time.Second),
		nc.ReconnectWait(1 * time.Second),
	}
	six_hour_before := time.Now().Add(-6 * time.Hour)

	subscribeOptions := []nc.SubOpt{
		nc.StartTime(six_hour_before),
		nc.AckExplicit(),
	}

	jsConfig := nats.JetStreamConfig{
		Disabled:         false,
		AutoProvision:    true,
		ConnectOptions:   nil,
		SubscribeOptions: subscribeOptions,
		PublishOptions:   nil,
		TrackMsgId:       false,
		AckAsync:         false,
		DurablePrefix:    "",
	}

	return nats.NewSubscriber(nats.SubscriberConfig{
		URL:              natsUrl,
		SubscribersCount: 1,
		NatsOptions:      options,
		JetStream:        jsConfig,
		Unmarshaler:      &nats.GobMarshaler{},
	}, watermill.NewStdLogger(false, false))
}

func CreateNATSPublisher(natsURL string) (*nats.Publisher, error) {
	if natsURL[:len(NatsURLPrefix)] != NatsURLPrefix {
		panic("natsUrl must start with " + NatsURLPrefix)
	}

	logger := watermill.NewStdLogger(false, false)
	marshaler := &nats.GobMarshaler{}
	options := []nc.Option{
		nc.RetryOnFailedConnect(true),
		nc.Timeout(30 * time.Second),
		nc.ReconnectWait(1 * time.Second),
	}
	subscribeOptions := []nc.SubOpt{
		nc.DeliverAll(),
		nc.AckExplicit(),
	}

	jsConfig := nats.JetStreamConfig{
		Disabled:         false,
		AutoProvision:    true,
		ConnectOptions:   nil,
		SubscribeOptions: subscribeOptions,
		PublishOptions:   nil,
		TrackMsgId:       false,
		AckAsync:         false,
		DurablePrefix:    "",
	}

	publisher, err := nats.NewPublisher(
		nats.PublisherConfig{
			URL:         natsURL,
			NatsOptions: options,
			Marshaler:   marshaler,
			JetStream:   jsConfig,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	return publisher, nil
}

func SubscribeToNATS(topic string) (<-chan *message.Message, error) {
	natsURL := utils.GetCredUnsafe("ATRO_NATS_URL")
	l := slog.Default().With("topic", topic, "nats_url", natsURL)

	l.Info("Creating NATS Subscriber")
	subscriber, err := CreateNATSSubscriber(natsURL)
	if err != nil {
		l.Error("Failed to create NATS Subscriber", "error", err)

		return nil, err
	}

	l.Info("Subscribing to NATS topic")
	messages, err := subscriber.Subscribe(context.Background(), topic)
	if err != nil {
		l.Error("Failed to subscribe to NATS topic", "error", err)

		return nil, err
	}

	l.Info("Successfully subscribed to NATS topic")

	return messages, nil
}

func PublishToNATS(topic string, event PubSubEvent) error {
	// Might be too "thick" on logging here, will reduce it later if needed

	natsUrl := utils.GetCredUnsafe("ATRO_NATS_URL")
	l := slog.Default().With("topic", topic, "nats_url", natsUrl)

	l.Info("Creating NATS Publisher")
	publisher, err := CreateNATSPublisher(natsUrl)
	if err != nil {
		l.Error("Failed to create NATS Publisher", "error", err)
		return err
	}

	l.Info("Marshalling event to bytes")
	eventBytes, err := json.Marshal(event)
	if err != nil {
		l.Error("Failed to marshal event", "error", err)
		return err
	}

	msgUUID := watermill.NewUUID()
	msg := message.NewMessage(msgUUID, message.Payload(eventBytes))
	l.Info("Created watermill message", "uuid", msgUUID)

	l.Info("Publishing message to NATS")
	if err := publisher.Publish("events", msg); err != nil {
		l.Error("Failed to publish event to NATS", "error", err)
		return err
	}

	l.Info("Succesfully published event to NATS")
	return nil
}
