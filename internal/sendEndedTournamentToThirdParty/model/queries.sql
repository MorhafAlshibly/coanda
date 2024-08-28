-- name: GetEndedTournaments :many
SELECT name,
    tournament_interval,
    tournament_started_at
FROM ranked_tournament
WHERE tournament_started_at < ?
    AND tournament_interval = ?
    AND (
        sent_to_third_party_at IS NULL
        OR sent_to_third_party_at > NOW()
    )
GROUP BY name,
    tournament_interval,
    tournament_started_at
ORDER BY name ASC,
    tournament_interval ASC,
    tournament_started_at DESC
LIMIT ? OFFSET ?;
-- name: GetEndedTournamentUsers :many
SELECT id,
    name,
    tournament_interval,
    user_id,
    score,
    ranking,
    data,
    tournament_started_at,
    created_at,
    updated_at
FROM ranked_tournament
WHERE name = ?
    AND tournament_started_at = ?
    AND tournament_interval = ?
    AND (
        sent_to_third_party_at IS NULL
        OR sent_to_third_party_at > NOW()
    )
    AND ranking <= ?;
-- name: UpdateTournamentSentToThirdParty :execresult
UPDATE tournament
SET sent_to_third_party_at = NOW()
WHERE name = ?
    AND tournament_started_at = ?
    AND tournament_interval = ?
    AND (
        sent_to_third_party_at IS NULL
        OR sent_to_third_party_at > NOW()
    );