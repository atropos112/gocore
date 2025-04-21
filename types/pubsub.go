package types

import (
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill"
)

// Event  represents an event that is sent over the pubsub system.
type Event struct {
	ID      string    `json:"id"`
	Source  string    `json:"source"`
	Type    string    `json:"type"`
	Data    []byte    `json:"data"`
	Time    time.Time `json:"time"`
	Subject string    `json:"subject"`
}

// MakeSimpleEvent makes a simple Event where the time is now, and ID is a UUID.
func MakeSimpleEvent(source, eventType, subject string, data any) (*Event, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	event := Event{
		ID:      watermill.NewUUID(),
		Source:  source,
		Type:    eventType,
		Data:    dataBytes,
		Time:    time.Now(),
		Subject: subject,
	}

	return &event, nil
}
