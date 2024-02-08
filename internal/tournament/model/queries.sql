-- name: GetTournament :one
SELECT name,
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
    AND tournament_interval = ?
    AND user_id = ?
    AND tournament_started_at = ?
LIMIT 1;
-- name: GetTournaments :many
SELECT name,
    tournament_interval,
    user_id,
    score,
    ranking,
    data,
    tournament_started_at,
    created_at,
    updated_at
FROM ranked_tournament
WHERE (
        name = sqlc.narg(name)
        OR user_id = sqlc.narg(user_id)
    )
    AND tournament_interval = ?
    AND tournament_started_at = ?
ORDER BY name ASC,
    tournament_interval ASC,
    score DESC
LIMIT ? OFFSET ?;
-- name: CreateTournament :execresult
INSERT INTO tournament (
        name,
        tournament_interval,
        user_id,
        score,
        data,
        tournament_started_at
    )
VALUES (?, ?, ?, ?, ?, ?);
-- name: DeleteTournament :execresult
DELETE FROM tournament
WHERE name = ?
    AND tournament_interval = ?
    AND user_id = ?
    AND tournament_started_at = ?
LIMIT 1;
-- name: UpdateTournamentScore :execresult
UPDATE tournament
SET score = CASE
        WHEN sqlc.arg(score) IS NOT NULL THEN sqlc.arg(score) + CASE
            WHEN CAST(sqlc.arg(increment_score) as unsigned) != 0 THEN score
            ELSE 0
        END
        ELSE score
    END
WHERE name = ?
    AND tournament_interval = ?
    AND user_id = ?
    AND tournament_started_at = ?
LIMIT 1;
-- name: UpdateTournamentData :execresult
UPDATE tournament
SET data = ?
WHERE name = ?
    AND tournament_interval = ?
    AND user_id = ?
    AND tournament_started_at = ?
LIMIT 1;
-- name: GetTournamentsBeforeWipe :many
SELECT name,
    tournament_interval,
    user_id,
    score,
    ranking,
    data,
    tournament_started_at,
    created_at,
    updated_at
FROM ranked_tournament
WHERE tournament_started_at < ?
    AND tournament_interval = ?
LIMIT ? OFFSET ?;
-- name: WipeTournaments :execresult
DELETE FROM tournament
WHERE tournament_started_at < ?
    AND tournament_interval = ?;