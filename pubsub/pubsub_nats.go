package pubsub

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	. "github.com/atropos112/gocore/logging"
	"github.com/atropos112/gocore/utils"
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
	}, InitWatermillLogger(),
	)
}

func CreateNATSPublisher(natsURL string) (*nats.Publisher, error) {
	if natsURL[:len(NatsURLPrefix)] != NatsURLPrefix {
		panic("natsUrl must start with " + NatsURLPrefix)
	}

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
		InitWatermillLogger(),
	)
	if err != nil {
		panic(err)
	}

	return publisher, nil
}

func SubscribeToNATS(topic string) (<-chan *message.Message, error) {
	natsURL, err := utils.GetCred("ATRO_NATS_URL")

	if _, ok := err.(*utils.NoCredFoundError); ok {
		slog.Default().Warn("Failed to get NATS URL, using default")
		natsURL = "nats://nats:4222" // Default NATS URL, the tailscale one.
	}

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

func PublishToNATS(topic string, event PublishableObject) error {
	// Might be too "thick" on logging here, will reduce it later if needed

	natsURL, err := utils.GetCred("ATRO_NATS_URL")

	if _, ok := err.(*utils.NoCredFoundError); ok {
		slog.Default().Warn("Failed to get NATS URL, setting for default")
		natsURL = "nats://nats:4222" // Default NATS URL, the tailscale one.
	}

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
	l.Info("Created watermill message", "uuid", msgUUID)

	l.Info("Publishing message to NATS")
	if err := publisher.Publish("events", msg); err != nil {
		l.Error("Failed to publish event to NATS", "error", err)
		return err
	}

	l.Info("Succesfully published event to NATS")
	return nil
}

func ErrorAlertAndDie(l *slog.Logger, source, msg string, args ...any) {
	ErrorAlert(l, source, msg, args...)
	panic("Failed with message: " + msg)
}

func ErrorAlert(l *slog.Logger, source, msg string, args ...any) {
	l.Error(msg, args...)

	argsMap := make(map[string]any)
	for i := 0; i < len(args); i += 2 {
		argsMap[args[i].(string)] = args[i+1]
	}

	PublishToNATS(NatsErrorsTopic, PubSubError{
		Source:  source,
		Message: msg,
		Args:    argsMap,
	})
}

func GetCredOrAlertAndDie(l *slog.Logger, source, value string) string {
	cred, err := utils.GetCred(value)
	if err != nil {
		ErrorAlertAndDie(l, source, "Failed to get credential", "error", err)
	}

	return cred
}

type structHandler struct {
	// we can add some dependencies here
}

func (h structHandler) Handler(msg *message.Message) error {
	// handle the message here
	return nil
}

func Test() {
	logger := InitWatermillLogger()
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	// Politely close the router when the program exits
	router.AddPlugin(plugin.SignalsHandler)

	router.AddMiddleware(
		// CorrelationID will copy the correlation id from the incoming message's metadata to the produced messages
		middleware.CorrelationID,

		// Recover from panics
		middleware.Recoverer,
	)

	subscriber, err := CreateNATSSubscriber(GetNATSURLFromEnv())
	if err != nil {
		panic(err)
	}

	router.AddNoPublisherHandler(
		"example_handler",
		NatsTestTopic, // Test topic name.
		subscriber,
		structHandler{}.Handler,
	)

	// handler := router.AddHandler(
	// 	"example_handler",
	// 	NatsTestTopic, // Test topic name.
	// 	subscriber,
	// 	NatsTestTopic, // Test topic name.
	// 	publisher,
	// )
	ctx := context.Background()

	if err := router.Run(ctx); err != nil {
		panic(err)
	}
}
