-- name: GetAgedMatchmakingTickets :many
SELECT id,
    matchmaking_match_id,
    elo_window,
    data,
    expires_at,
    created_at,
    updated_at
FROM matchmaking_ticket
WHERE expires_at < NOW()
    AND matchmaking_match_id IS NULL
    AND elo_window >= sqlc.arg(elo_window_max)
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
WHERE mt.id = ?
GROUP BY ma.id,
    ma.name
ORDER BY ticket_count DESC
LIMIT 1;
-- name: AddMatchIDToTicket :execresult
UPDATE matchmaking_ticket
SET matchmaking_match_id = ?
WHERE id = ?
    AND matchmaking_match_id IS NULL;
-- name: IncrementEloWindow :execresult
UPDATE matchmaking_ticket
SET elo_window = elo_window + sqlc.arg(elo_window_increment)
WHERE expires_at < NOW()
    AND matchmaking_match_id IS NULL
    AND elo_window < sqlc.arg(elo_window_max);
-- name: GetNonAgedMatchmakingTickets :many
SELECT id,
    matchmaking_match_id,
    elo_window,
    data,
    expires_at,
    created_at,
    updated_at
FROM matchmaking_ticket
WHERE expires_at < NOW()
    AND matchmaking_match_id IS NULL
    AND elo_window < sqlc.arg(elo_window_max)
LIMIT ? OFFSET ?;
-- name: GetClosestMatch :one
WITH ticket_elo AS (
    SELECT AVG(mu.elo) - mt.elo_window AS min_elo,
        AVG(mu.elo) + mt.elo_window AS max_elo,
        COUNT(mu.id) AS player_count
    FROM matchmaking_ticket mt
        JOIN matchmaking_ticket_user mtu ON mt.id = mtu.matchmaking_ticket_id
        JOIN matchmaking_user mu ON mtu.matchmaking_user_id = mu.id
    WHERE mt.id = sqlc.arg(ticket_id)
),
match_elo AS (
    SELECT mm.id AS match_id,
        mm.matchmaking_arena_id,
        AVG(mu.elo) AS avg_elo,
        COUNT(mu.id) AS current_player_count
    FROM matchmaking_match mm
        JOIN matchmaking_ticket mt ON mm.id = mt.matchmaking_match_id
        JOIN matchmaking_ticket_user mtu ON mt.id = mtu.matchmaking_ticket_id
        JOIN matchmaking_user mu ON mtu.matchmaking_user_id = mu.id
        JOIN ticket_elo te ON mm.matchmaking_arena_id = te.matchmaking_arena_id
    WHERE mm.locked_at < NOW()
    GROUP BY mm.id,
        mm.matchmaking_arena_id
    HAVING avg_elo BETWEEN MIN(te.min_elo) AND MAX(te.max_elo)
        AND te.player_count <= (
            SELECT max_players
            FROM matchmaking_arena
            WHERE id = mm.matchmaking_arena_id
        ) - current_player_count
)
SELECT mm.id,
    mm.matchmaking_arena_id,
    mm.data,
    mm.locked_at,
    mm.started_at,
    mm.ended_at,
    mm.created_at,
    mm.updated_at
FROM matchmaking_match mm
    JOIN match_elo me ON mm.id = me.match_id
ORDER BY ABS(
        me.avg_elo - (
            SELECT AVG(mu.elo)
            FROM matchmaking_ticket mt
                JOIN matchmaking_ticket_user mtu ON mt.id = mtu.matchmaking_ticket_id
                JOIN matchmaking_user mu ON mtu.matchmaking_user_id = mu.id
            WHERE mt.id = sqlc.arg(ticket_id)
        )
    )
LIMIT 1;