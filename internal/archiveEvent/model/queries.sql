-- name: GetEndedEventRoundUsers :many
SELECT eru.id,
    eru.event_user_id,
    eru.event_round_id,
    eru.result,
    eru.data,
    eru.created_at,
    eru.updated_at
FROM event_round_user eru
WHERE (
        SELECT er.ended_at
        FROM event_round er
        WHERE er.id = eru.event_round_id
    ) < NOW()
ORDER BY eru.id ASC
LIMIT ? OFFSET ?;
-- name: DeleteEndedEventRoundUsers :execresult
DELETE FROM event_round_user eru
WHERE eru.event_round_id IN (
        SELECT er.id
        FROM event_round er
        WHERE er.ended_at < NOW()
    )
    AND eru.id >= sqlc.arg(min_id)
    AND eru.id <= sqlc.arg(max_id);
-- name: GetEndedEventRounds :many
SELECT id,
    event_id,
    name,
    scoring,
    data,
    ended_at,
    created_at,
    updated_at
FROM event_round
WHERE ended_at < NOW()
ORDER BY id ASC
LIMIT ? OFFSET ?;
-- name: DeleteEndedEventRounds :execresult
DELETE FROM event_round
WHERE ended_at < NOW()
    AND id >= sqlc.arg(min_id)
    AND id <= sqlc.arg(max_id);
-- name: GetEndedEventUsers :many
SELECT eu.id,
    eu.event_id,
    eu.user_id,
    eu.data,
    eu.created_at,
    eu.updated_at
FROM event_user eu
WHERE (
        SELECT er.ended_at
        FROM event_round er
        WHERE er.event_id = event_user.event_id
        ORDER BY er.ended_at DESC
        LIMIT 1
    ) < NOW()
ORDER BY eu.id ASC
LIMIT ? OFFSET ?;
-- name: DeleteEndedEventUsers :execresult
DELETE FROM event_user eu
WHERE (
        SELECT er.ended_at
        FROM event_round er
        WHERE er.event_id = event_user.event_id
        ORDER BY er.ended_at DESC
        LIMIT 1
    ) < NOW()
    AND eu.id >= sqlc.arg(min_id)
    AND eu.id <= sqlc.arg(max_id);
-- name: GetEndedEvents :many
SELECT e.id,
    e.name,
    e.data,
    e.started_at,
    e.created_at,
    e.updated_at
FROM event e
WHERE (
        SELECT er.ended_at
        FROM event_round er
        WHERE er.event_id = event.id
        ORDER BY er.ended_at DESC
        LIMIT 1
    ) < NOW()
ORDER BY e.id ASC
LIMIT ? OFFSET ?;
-- name: DeleteEndedEvents :execresult
DELETE FROM event e
WHERE (
        SELECT er.ended_at
        FROM event_round er
        WHERE er.event_id = event.id
        ORDER BY er.ended_at DESC
        LIMIT 1
    ) < NOW()
    AND e.id >= sqlc.arg(min_id)
    AND e.id <= sqlc.arg(max_id);