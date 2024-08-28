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
    AND er.ended_at < NOW()
    AND (
        er.sent_to_third_party_at IS NULL
        OR er.sent_to_third_party_at > NOW()
    )
ORDER BY erl.ranking ASC
LIMIT ?;
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