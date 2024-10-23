// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package events

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const eventExists = `-- name: EventExists :one
SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)
`

func (q *Queries) EventExists(ctx context.Context, id string) (bool, error) {
	row := q.db.QueryRow(ctx, eventExists, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getEvent = `-- name: GetEvent :one
SELECT id, source, type, specversion, datacontenttype, data, time, subject FROM events WHERE id = $1
`

func (q *Queries) GetEvent(ctx context.Context, id string) (Event, error) {
	row := q.db.QueryRow(ctx, getEvent, id)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Source,
		&i.Type,
		&i.Specversion,
		&i.Datacontenttype,
		&i.Data,
		&i.Time,
		&i.Subject,
	)
	return i, err
}

const getEventsUpTo = `-- name: GetEventsUpTo :many
SELECT id, source, type, specversion, datacontenttype, data, time, subject FROM events WHERE time > $1
`

func (q *Queries) GetEventsUpTo(ctx context.Context, time pgtype.Timestamptz) ([]Event, error) {
	rows, err := q.db.Query(ctx, getEventsUpTo, time)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.Source,
			&i.Type,
			&i.Specversion,
			&i.Datacontenttype,
			&i.Data,
			&i.Time,
			&i.Subject,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertEvent = `-- name: InsertEvent :one
INSERT INTO events (id, source, type, specversion, datacontenttype, data, time, subject) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (id) DO NOTHING RETURNING id, source, type, specversion, datacontenttype, data, time, subject
`

type InsertEventParams struct {
	ID              string
	Source          string
	Type            string
	Specversion     string
	Datacontenttype pgtype.Text
	Data            []byte
	Time            pgtype.Timestamptz
	Subject         pgtype.Text
}

func (q *Queries) InsertEvent(ctx context.Context, arg InsertEventParams) (Event, error) {
	row := q.db.QueryRow(ctx, insertEvent,
		arg.ID,
		arg.Source,
		arg.Type,
		arg.Specversion,
		arg.Datacontenttype,
		arg.Data,
		arg.Time,
		arg.Subject,
	)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Source,
		&i.Type,
		&i.Specversion,
		&i.Datacontenttype,
		&i.Data,
		&i.Time,
		&i.Subject,
	)
	return i, err
}
