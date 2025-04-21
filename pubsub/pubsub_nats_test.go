// nolint
package pubsub_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/atropos112/atrogolib/pubsub"
	"github.com/atropos112/atrogolib/types"
)

func TestRoundTrip(t *testing.T) {
	p, err := pubsub.MakeNATSPublisher()
	if err != nil {
		t.Fatalf("Failed to make NATS publisher: %v", err)
	}

	event, err := types.MakeSimpleEvent("test_atrogolib", "test", "test", map[string]any{"a": 1, "b": 1})
	if err != nil {
		t.Fatalf("Failed to make event: %v", err)
	}

	now := time.Now().UTC()

	err = p.Publish("test", *event, nil)
	if err != nil {
		t.Fatalf("Failed to publish event: %v", err)
	}

	sub, err := pubsub.MakeNATSSubscriber(now)
	if err != nil {
		t.Fatalf("Failed to make NATS subscriber: %v", err)
	}

	eventChan, err := sub.Subscribe("test")
	if err != nil {
		t.Fatalf("Failed to subscribe to event: %v", err)
	}

	select {
	case e := <-eventChan:
		if event.ID == e.ID {
			if event.Source != e.Source {
				t.Errorf("Expected source %s, got %s", event.Source, e.Source)
			}
			if event.Type != e.Type {
				t.Errorf("Expected type %s, got %s", event.Type, e.Type)
			}
			if event.Subject != e.Subject {
				t.Errorf("Expected subject %s, got %s", event.Subject, e.Subject)
			}
			eventData := map[string]any{}
			if er := json.Unmarshal(e.Data, &eventData); er != nil {
				t.Fatalf("Failed to unmarshal event data: %v", er)
			}

			eData := map[string]any{}
			if er := json.Unmarshal(event.Data, &eData); er != nil {
				t.Fatalf("Failed to unmarshal event data: %v", er)
			}

			for k, v := range eventData {
				if eData[k] != v {
					t.Errorf("Expected data %s: %v, got %v", k, v, eData[k])
				}
			}

			break
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("Timed out waiting for event")
	}
}
