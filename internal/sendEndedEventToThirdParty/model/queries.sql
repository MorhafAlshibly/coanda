-- name: GetEndedEventRounds :many
SELECt id,
    event_id,
    name,
    scoring,
    data,
    ended_at,
    created_at,
    updated_at
FROM event_round
WHERE ended_at < NOW()
    AND (
        sent_to_third_party_at IS NULL
        OR sent_to_third_party_at > NOW()
    )
ORDER BY id ASC
LIMIT ? OFFSET ?;
-- name: GetEndedEventRoundLeaderboard :many
SELECT erl.id,
    erl.event_id,
    erl.round_name,
    erl.event_user_id,
    erl.client_user_id,
    erl.event_round_id,
    erl.result,
    erl.score,
    erl.ranking,
    erl.data,
    erl.created_at,
    erl.updated_at
FROM event_round_leaderboard erl
    JOIN event_round er ON erl.event_round_id = er.id
WHERE erl.event_round_id = ?
    AND erl.ranking <= ?
    AND er.ended_at < NOW()
    AND (
        er.sent_to_third_party_at IS NULL
        OR er.sent_to_third_party_at > NOW()
    )
ORDER BY erl.ranking ASC;
-- name: UpdateEventRoundSentToThirdParty :execresult
UPDATE event_round
SET sent_to_third_party_at = NOW()
WHERE id = ?
    AND ended_at < NOW()
    AND (
        sent_to_third_party_at IS NULL
        OR sent_to_third_party_at > NOW()
    )
LIMIT 1;
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
        WHERE er.event_id = e.id
        ORDER BY er.ended_at DESC
        LIMIT 1
    ) < NOW()
    AND (
        e.sent_to_third_party_at IS NULL
        OR e.sent_to_third_party_at > NOW()
    )
ORDER BY e.id ASC
LIMIT ? OFFSET ?;
-- name: GetEndedEventLeaderboard :many
SELECT el.id,
    el.event_id,
    el.client_user_id,
    el.score,
    el.ranking,
    el.data,
    el.created_at,
    el.updated_at
FROM event_leaderboard el
    JOIN event e ON el.event_id = e.id
WHERE e.id = ?
    AND el.ranking <= ?
    AND (
        SELECT er.ended_at
        FROM event_round er
        WHERE er.event_id = e.id
        ORDER BY er.ended_at DESC
        LIMIT 1
    ) < NOW()
    AND (
        e.sent_to_third_party_at IS NULL
        OR e.sent_to_third_party_at > NOW()
    )
ORDER BY el.ranking ASC;
-- name: UpdateEventSentToThirdParty :execresult
UPDATE event e
SET e.sent_to_third_party_at = NOW()
WHERE e.id = ?
    AND (
        SELECT er.ended_at
        FROM event_round er
        WHERE er.event_id = e.id
        ORDER BY er.ended_at DESC
        LIMIT 1
    ) < NOW()
    AND (
        e.sent_to_third_party_at IS NULL
        OR e.sent_to_third_party_at > NOW()
    )
LIMIT 1;