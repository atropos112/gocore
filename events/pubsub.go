package events

import (
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// PubSubEvent represents an event that is sent over the pubsub system.
type PubSubEvent struct {
	ID              string    `json:"id"`
	Source          string    `json:"source"`
	Type            string    `json:"type"`
	Specversion     string    `json:"specversion"`
	Datacontenttype string    `json:"datacontenttype"`
	Data            []byte    `json:"data"`
	Time            time.Time `json:"time"`
	Subject         string    `json:"subject"`
}

// ToDBEvent transforms PubSubEvent to DBEvent.
func (e *PubSubEvent) ToDBEvent() Event {
	return Event{
		ID:              e.ID,
		Source:          e.Source,
		Type:            e.Type,
		Specversion:     e.Specversion,
		Datacontenttype: pgtype.Text{String: e.Datacontenttype},
		Data:            e.Data,
		Time:            pgtype.Timestamptz{Time: e.Time},
		Subject:         pgtype.Text{String: e.Subject},
	}
}

// ToPubSubEvent transforms DBEvent to PubSubEvent.
func (e *Event) ToPubSubEvent() PubSubEvent {
	return PubSubEvent{
		ID:              e.ID,
		Source:          e.Source,
		Type:            e.Type,
		Specversion:     e.Specversion,
		Datacontenttype: e.Datacontenttype.String,
		Data:            e.Data,
		Time:            e.Time.Time,
		Subject:         e.Subject.String,
	}
}

// GetData provdies a way to get the data from the event.
func (e *PubSubEvent) GetData(v *any) error {
	return json.Unmarshal(e.Data, v)
}
