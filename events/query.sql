-- name: InsertEvent :one
INSERT INTO events (id, source, type, specversion, datacontenttype, data, time, subject) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetEvent :one
SELECT * FROM events WHERE id = $1;

-- name: GetEventsUpTo :many
SELECT * FROM events WHERE time > $1;

-- name: EventExists :one
SELECT EXISTS(SELECT 1 FROM events WHERE id = $1);
