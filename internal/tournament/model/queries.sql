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
-- name: GetTournamentsBeforeWipe :many
SELECT *
FROM ranked_tournament
WHERE tournament_started_at < ?
    AND tournament_interval = ?
LIMIT ? OFFSET ?;
-- name: WipeTournaments :execresult
DELETE FROM tournament
WHERE tournament_started_at < ?
    AND tournament_interval = ?;