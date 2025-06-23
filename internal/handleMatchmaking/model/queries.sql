-- name: GetAgedMatchmakingTickets :many
SELECT id,
    matchmaking_match_id,
    data,
    created_at,
    updated_at
FROM matchmaking_ticket
WHERE matchmaking_match_id IS NULL -- Ensure difference in seconds between now and crearted at multiplied by elo_window_increment_per_second is greater than elo_window_max
    AND (
        TIMESTAMPDIFF(SECOND, created_at, NOW()) * CAST(
            sqlc.arg(elo_window_increment_per_second) AS UNSIGNED INTEGER
        )
    ) > CAST(sqlc.arg(elo_window_max) AS UNSIGNED INTEGER)
LIMIT ? OFFSET ?;
-- name: CreateMatch :execresult
INSERT INTO matchmaking_match (matchmaking_arena_id, data)
VALUES (?, "{}");
-- name: GetMostPopularArenaOnTicket :one
SELECT ma.id,
    ma.name,
    ma.min_players,
    ma.max_players_per_ticket,
    ma.max_players,
    ma.data,
    ma.created_at,
    ma.updated_at,
    COUNT(mt.id) AS ticket_count
FROM matchmaking_ticket mt
    JOIN matchmaking_ticket_arena mta ON mt.id = mta.matchmaking_ticket_id
    JOIN matchmaking_arena ma ON mta.matchmaking_arena_id = ma.id
WHERE ma.id IN (
        SELECT mta.matchmaking_arena_id
        FROM matchmaking_ticket_arena mta
        WHERE mta.matchmaking_ticket_id = ?
    )
GROUP BY ma.id
ORDER BY ticket_count DESC
LIMIT 1;
-- name: AddMatchIDToTicket :execresult
UPDATE matchmaking_ticket
SET matchmaking_match_id = ?
WHERE id = ?;
-- name: GetNonAgedMatchmakingTickets :many
SELECT id,
    matchmaking_match_id,
    data,
    created_at,
    updated_at
FROM matchmaking_ticket
WHERE matchmaking_match_id IS NULL
    AND (
        TIMESTAMPDIFF(SECOND, created_at, NOW()) * CAST(
            sqlc.arg(elo_window_increment_per_second) AS UNSIGNED INTEGER
        )
    ) <= CAST(sqlc.arg(elo_window_max) AS UNSIGNED INTEGER)
LIMIT ? OFFSET ?;
-- name: LockTicketForUpdate :one
SELECT id,
    matchmaking_match_id,
    data,
    created_at,
    updated_at
FROM matchmaking_ticket
WHERE id = ?
    AND matchmaking_match_id IS NULL
LIMIT 1 FOR
UPDATE;
-- name: GetClosestMatch :one
WITH ticket_info AS (
    SELECT mu.matchmaking_ticket_id,
        COUNT(DISTINCT mu.id) AS user_count,
        AVG(mu.elo) AS avg_elo
    FROM matchmaking_user mu
    WHERE mu.matchmaking_ticket_id IS NOT NULL
        AND mu.matchmaking_ticket_id = sqlc.narg(ticket_id)
    GROUP BY mu.matchmaking_ticket_id
),
ticket_arenas AS (
    SELECT mta.matchmaking_arena_id
    FROM matchmaking_ticket_arena mta
    WHERE mta.matchmaking_ticket_id = sqlc.narg(ticket_id)
)
SELECT mm.id AS match_id,
    ma.name AS arena_name,
    ma.max_players,
    current_players,
    (ma.max_players - current_players) AS remaining_capacity,
    ti.user_count AS ticket_user_count,
    CAST(match_avg_elo AS UNSIGNED INTEGER) AS match_avg_elo,
    CAST(ti.avg_elo AS UNSIGNED INTEGER) AS ticket_avg_elo,
    CAST(
        ABS(match_avg_elo - ti.avg_elo) AS UNSIGNED INTEGER
    ) AS elo_difference,
    mm.locked_at -- Added for visibility
FROM matchmaking_match mm
    JOIN matchmaking_arena ma ON mm.matchmaking_arena_id = ma.id
    JOIN ticket_arenas ta ON mm.matchmaking_arena_id = ta.matchmaking_arena_id
    JOIN (
        SELECT mt.matchmaking_match_id,
            COUNT(DISTINCT mu.id) AS current_players,
            AVG(mu.elo) AS match_avg_elo
        FROM matchmaking_ticket mt
            JOIN matchmaking_user mu ON mu.matchmaking_ticket_id = mt.id
        WHERE mt.matchmaking_match_id IS NOT NULL
        GROUP BY mt.matchmaking_match_id
    ) match_stats ON mm.id = match_stats.matchmaking_match_id
    JOIN ticket_info ti
WHERE ti.user_count <= (ma.max_players - current_players)
    AND (
        mm.locked_at IS NULL
        OR mm.locked_at > NOW()
    )
ORDER BY elo_difference ASC
LIMIT 1;