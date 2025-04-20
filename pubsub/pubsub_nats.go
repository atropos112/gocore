package pubsub

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	. "github.com/atropos112/atrogolib/logging"
	. "github.com/atropos112/atrogolib/types"
	"github.com/atropos112/atrogolib/utils"
	nc "github.com/nats-io/nats.go"
)

// NatsUrlPrefix  is the prefix for NATS URLs.
const NatsURLPrefix = "nats://"

func GetNATSURLFromEnv() string {
	natsURL, err := utils.GetCred("ATRO_NATS_URL")

	if _, ok := err.(*utils.NoCredFoundError); ok {
		slog.Default().Warn("Failed to get NATS URL, using default")
		natsURL = "nats://nats:4222" // Default NATS URL, the tailscale one.
	}

	return natsURL
}

// CreateSubscriber creates a new NATS publisher.
// Can then use it to subscribe to topic like
// messages, err := subscriber.Subscribe(context.Background(), "example.topic")
func CreateNATSSubscriber(natsUrl string, startTime time.Time) (*nats.Subscriber, error) {
	if natsUrl[:len(NatsURLPrefix)] != NatsURLPrefix {
		panic("natsUrl must start with " + NatsURLPrefix)
	}
	marshaler := &nats.JSONMarshaler{}

	options := []nc.Option{
		nc.RetryOnFailedConnect(true),
		nc.Timeout(30 * time.Second),
		nc.ReconnectWait(1 * time.Second),
	}

	subscribeOptions := []nc.SubOpt{
		nc.StartTime(startTime),
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
		Unmarshaler:      marshaler,
	}, InitWatermillLogger(),
	)
}

func CreateNATSPublisher(natsURL string) (*nats.Publisher, error) {
	if natsURL[:len(NatsURLPrefix)] != NatsURLPrefix {
		panic("natsUrl must start with " + NatsURLPrefix)
	}

	marshaler := &nats.JSONMarshaler{}
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
		InitWatermillLogger(),
	)
	if err != nil {
		panic(err)
	}

	return publisher, nil
}

func PublishToNATS(topic string, event PublishableObject, metadata map[string]string) error {
	// Might be too "thick" on logging here, will reduce it later if needed
	natsURL := utils.GetCredUnsafe("ATRO_NATS_URL")

	l := slog.Default().With("topic", topic, "nats_url", natsURL)

	l.Info("Creating NATS Publisher")
	publisher, err := CreateNATSPublisher(natsURL)
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
	msg.Metadata = metadata
	l.Info("Created watermill message", "uuid", msgUUID)

	l.Info("Publishing message to NATS")
	if err := publisher.Publish(topic, msg); err != nil {
		l.Error("Failed to publish event to NATS", "error", err)
		return err
	}

	l.Info("Succesfully published event to NATS")
	return nil
}

func NewNATSAlertContext(l *slog.Logger, source string) *AlertContext {
	return &AlertContext{
		Source: source,
		Logger: l,
		Publish: func(obj PublishableObject, source string) {
			err := PublishToNATS(NatsErrorsTopic, obj, map[string]string{"source": source})
			if err != nil {
				l.Error(err.Error())
				panic("Failed to publish error to NATS")
			}
		},
	}
}

func NewEventNATSRouter(l *slog.Logger) *message.Router {
	watermillLogger := watermill.NewSlogLogger(l)

	router, err := message.NewRouter(message.RouterConfig{}, watermillLogger)
	if err != nil {
		panic(err)
	}

	router.AddPlugin(plugin.SignalsHandler)

	router.AddMiddleware(
		// Add timeout to context, in case of a timeout, the message will be nacked.
		middleware.Timeout(time.Second*10),

		// Add correlation ID to context,
		middleware.CorrelationID,
	)

	return router
}
