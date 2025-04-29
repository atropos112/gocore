package pubsub

import (
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/atropos112/atrogolib/logging"
	"github.com/atropos112/atrogolib/types"
	"github.com/atropos112/atrogolib/utils"
	nc "github.com/nats-io/nats.go"
)

type NATSPublisher struct {
	Publisher *nats.Publisher
	NatsURL   string
	l         *slog.Logger
}

// Purely to cause bang when we fail to satisty the interface.
var _ Publisher = &NATSPublisher{} //nolint:exhaustruct // Just to test we satistfy the interface

func (p *NATSPublisher) Publish(topic string, event types.Event, metadata *map[string]string) error {
	msgUUID := watermill.NewUUID()
	msg := message.NewMessage(
		msgUUID,
		message.Payload(event.Data),
	)
	if metadata != nil {
		msg.Metadata = *metadata
	} else {
		msg.Metadata = make(map[string]string)
	}

	msg.Metadata["source"] = event.Source
	msg.Metadata["time"] = event.Time.UTC().Format(time.RFC3339Nano)
	msg.Metadata["type"] = event.Type

	if event.Subject != "" {
		msg.Metadata["subject"] = event.Subject
	}

	if err := p.Publisher.Publish(topic, msg); err != nil {
		return err
	}

	p.l.Info("Published message to NATS", "topic", topic, "UUID", msgUUID, "source", event.Source, "time", msg.Metadata["time"], "type", event.Type, "subject", event.Subject)

	return nil
}

// NatsURLPrefix is the prefix for NATS URLs.

func MakeNATSPublisher() (*NATSPublisher, error) {
	natsURL, err := utils.GetCred("ATRO_NATS_URL")
	if err != nil {
		return nil, err
	}

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
		logging.InitWatermillLogger(),
	)
	if err != nil {
		return nil, err
	}

	natsPublisher := NATSPublisher{
		Publisher: publisher,
		NatsURL:   natsURL,
		l:         slog.Default().With("nats_url", natsURL),
	}

	return &natsPublisher, nil
}
