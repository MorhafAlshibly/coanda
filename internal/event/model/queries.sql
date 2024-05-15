-- name: CreateEvent :execresult
INSERT INTO event (name, data, started_at)
VALUES (?, ?, ?);