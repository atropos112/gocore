package pubsub

import (
	"context"
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	. "github.com/atropos112/atrogolib/logging"
	. "github.com/atropos112/atrogolib/types"
	"github.com/atropos112/atrogolib/utils"
	nc "github.com/nats-io/nats.go"
	"github.com/negrel/assert"
)

type NATSSubscriber struct {
	Subscriber *nats.Subscriber
	NatsURL    string
	l          *slog.Logger
}

// Purely to cause bang when we fail to satisty the interface.
var _ Subscriber = &NATSSubscriber{} //nolint:exhaustruct // Just to test we satistfy the interface

func (s NATSSubscriber) Subscribe(topic string) (<-chan Event, error) {
	msgChan, err := s.Subscriber.Subscribe(context.Background(), topic)
	if err != nil {
		return nil, err
	}

	eventChan := make(chan Event)

	go func() {
		for msg := range msgChan {
			meta := msg.Metadata

			parsedTime, er := time.Parse(time.RFC3339Nano, meta["time"])
			assert.NoErrorf(er, "Tried to parse message form topic and failed at parsing time")
			if er != nil {
				continue
			}

			event := Event{
				ID:      msg.UUID,
				Source:  meta["source"],
				Type:    meta["type"],
				Time:    parsedTime,
				Subject: meta.Get("subject"),
				Data:    msg.Payload,
			}

			eventChan <- event

			msg.Ack()
		}
		close(eventChan) // Close the event channel when msgChan is closed
	}()

	return eventChan, nil
}

// MakeNATSSubscriber creates a new NATS subscriber.
func MakeNATSSubscriber(startTime time.Time) (*NATSSubscriber, error) {
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

	natsConfig := nats.SubscriberConfig{
		URL:              natsURL,
		SubscribersCount: 1,
		NatsOptions:      options,
		JetStream:        jsConfig,
		Unmarshaler:      marshaler,
	}

	subscriber, err := nats.NewSubscriber(natsConfig, InitWatermillLogger())
	if err != nil {
		return nil, err
	}

	natsSub := NATSSubscriber{
		Subscriber: subscriber,
		NatsURL:    natsURL,
		l:          slog.Default().With("nats_url", natsURL),
	}

	return &natsSub, nil
}
