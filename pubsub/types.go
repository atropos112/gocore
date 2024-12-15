package pubsub

import (
	"encoding/json"
	"time"

	"github.com/atropos112/gocore/events"
	"github.com/jackc/pgx/v5/pgtype"
)

type PublishableObject any

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

type PubSubError struct {
	Source  string         `json:"source"`
	Message string         `json:"message"`
	Args    map[string]any `json:"args"`
}

// ToDBEvent transforms PubSubEvent to DBEvent.
func (e *PubSubEvent) ToDBEvent() events.Event {
	return events.Event{
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

// ToDBInstertEventParams provides a way to insert the event into the database.
func (e *PubSubEvent) ToDBInstertEventParams() events.InsertEventParams {
	return events.InsertEventParams{
		ID:              e.ID,
		Source:          e.Source,
		Type:            e.Type,
		Specversion:     e.Specversion,
		Datacontenttype: pgtype.Text{String: e.Datacontenttype, Valid: true},
		Data:            e.Data,
		Time:            pgtype.Timestamptz{Time: e.Time, Valid: true},
		Subject:         pgtype.Text{String: e.Subject, Valid: true},
	}
}

// ToPubSubEvent transforms DBEvent to PubSubEvent.
func ToPubSubEvent(e *events.Event) PubSubEvent {
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
