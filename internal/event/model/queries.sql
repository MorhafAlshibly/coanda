-- name: CreateEvent :execresult
INSERT INTO event (name, data, started_at)
VALUES (?, ?, ?);
-- name: CreateEventRound :execresult
INSERT INTO event_round (event_id, name, data, scoring, ended_at)
VALUES (?, ?, ?, ?, ?);