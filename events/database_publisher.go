package events

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/atropos112/gocore/types"
	"github.com/atropos112/gocore/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const pSQLPort = 5432

type DatabasePublisher struct {
	Queries *Queries
	Context *context.Context
	Logger  *slog.Logger
}

func NewDatabasePublisher() DatabasePublisher {
	l := slog.Default()
	ctx := context.Background()

	l.Info("Sourcing credentials")

	dbpassword := utils.GetCredUnsafe("GOCORE_EVENTS_DB_PASSWORD")

	// On cluster host == psql-events-rw.events, locally host == events (tailscale)
	host := utils.GetCredUnsafe("GOCORE_EVENTS_DB_HOST")
	connString := "postgres://app:" + dbpassword + "@" + host + ":" + strconv.Itoa(pSQLPort) + "/app"

	l.Info("Connecting to PSQL", "host", host, "port", pSQLPort)
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		panic(err)
	}
	queries := New(conn)

	l.Info("Connected to PSQL", "host", host, "port", pSQLPort)
	return DatabasePublisher{
		Queries: queries,
		Context: &ctx,
		Logger:  l,
	}
}

// ToDBEvent transforms PubSubEvent to DBEvent.
func PublishableEventToDBEvent(e types.PublishableEvent) Event {
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

// ToDBInstertEventParams provides a way to insert the event into the database.
func ToDBInstertEventParams(e types.PublishableEvent) InsertEventParams {
	return InsertEventParams{
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
