package types

import (
	"encoding/json"
	"time"
)

type PublishableObject any

// PublishableEvent  represents an event that is sent over the pubsub system.
type PublishableEvent struct {
	ID              string    `json:"id"`
	Source          string    `json:"source"`
	Type            string    `json:"type"`
	Specversion     string    `json:"specversion"`
	Datacontenttype string    `json:"datacontenttype"`
	Data            []byte    `json:"data"`
	Time            time.Time `json:"time"`
	Subject         string    `json:"subject"`
}

// GetData provdies a way to get the data from the event.
func (e *PublishableEvent) GetData(v *any) error {
	return json.Unmarshal(e.Data, v)
}
