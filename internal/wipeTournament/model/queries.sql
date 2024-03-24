-- name: WipeTournaments :execresult
DELETE FROM tournament
WHERE tournament_started_at < ?
    AND tournament_interval = ?;
-- name: ArchiveTournaments :execresult
INSERT INTO archived_tournament (
        id,
        name,
        tournament_interval,
        user_id,
        score,
        data,
        tournament_started_at,
        created_at,
        updated_at
    )
SELECT t.id,
    t.name,
    t.tournament_interval,
    t.user_id,
    t.score,
    t.data,
    t.tournament_started_at,
    t.created_at,
    t.updated_at
FROM tournament t
WHERE t.tournament_started_at < ?
    AND t.tournament_interval = ?;