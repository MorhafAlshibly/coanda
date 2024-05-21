-- name: CreateEvent :execresult
INSERT INTO event (name, data, started_at)
VALUES (?, ?, ?);
-- name: CreateEventRound :execresult
INSERT INTO event_round (event_id, name, data, scoring, ended_at)
VALUES (?, ?, ?, ?, ?);
-- name: CreateEventUser :execresult
INSERT INTO event_user (event_id, user_id, data)
VALUES (?, ?, ?);
-- name: CreateEventRoundUser :execresult
INSERT INTO event_round_user (event_user_id, event_round_id, result, data)
SELECT ?,
    er.id,
    ?,
    ?
FROM event_round er
WHERE er.ended_at = (
        SELECT MIN(ended_at)
        FROM event_round
        WHERE ended_at > NOW()
    )
LIMIT 1;
-- name: GetEventByName :one
SELECT id,
    name,
    data,
    started_at,
    created_at,
    updated_at
FROM event
WHERE name = ?
LIMIT 1;
-- name: GetEventUserByEventIdAndUserId :one
SELECT id,
    event_id,
    user_id,
    data,
    created_at,
    updated_at
FROM event_user
WHERE event_id = ?
    AND user_id = ?
LIMIT 1;