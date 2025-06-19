-- name: CreateEvent :execresult
INSERT INTO event (name, data, started_at)
VALUES (?, ?, ?);
-- name: CreateEventRound :execresult
INSERT INTO event_round (event_id, name, data, scoring, ended_at)
VALUES (?, ?, ?, ?, ?);
-- name: CreateEventUser :execresult
INSERT INTO event_user (event_id, client_user_id, data)
VALUES (?, ?, ?);
-- name: CreateEventRoundUser :execresult
INSERT INTO event_round_user (event_user_id, event_round_id, result, data)
VALUES (?, ?, ?, ?);
-- name: UpdateEventRoundUser :execresult
UPDATE event_round_user eru
SET eru.result = ?,
    eru.data = ?
WHERE eru.event_user_id = ?
    AND eru.event_round_id = ?
LIMIT 1;
-- TODO: Revaluate if we should keep this query below
-- name: GetEventRoundUserByEventUserId :one
SELECT eru.event_user_id,
    eru.event_round_id,
    eru.result,
    eru.data,
    eru.created_at,
    eru.updated_at
FROM event_round_user eru
    JOIN event_round er ON eru.event_round_id = er.id
WHERE eru.event_user_id = ?
    AND er.ended_at > NOW()
ORDER BY er.ended_at ASC
LIMIT 1;
-- name: DeleteEventRoundUser :execresult
DELETE FROM event_round_user
WHERE id = ?
LIMIT 1;