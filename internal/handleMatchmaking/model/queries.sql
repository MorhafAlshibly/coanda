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
    AND elo_window > ?
LIMIT ? OFFSET ?;
-- name: CreateMatch :execresult
INSERT INTO matchmaking_match (matchmaking_arena_id, data)
VALUES (?, ?);
-- name: AddMatchIDToTicket :exec
UPDATE matchmaking_ticket
SET matchmaking_match_id = ?
WHERE id = ?
    AND matchmaking_match_id IS NULL;
-- name: ExpandEloWindow :exec
UPDATE matchmaking_ticket
SET elo_window = elo_window + ?
WHERE expires_at < NOW()
    AND matchmaking_match_id IS NULL
    AND elo_window < ?;
-- name: GetClosestMatch :one
WITH TicketElo AS (
    SELECT AVG(mue.elo) - mt.elo_window AS min_elo,
        AVG(mue.elo) + mt.elo_window AS max_elo,
        COUNT(mu.id) AS player_count,
        mue.matchmaking_arena_id
    FROM matchmaking_ticket mt
        JOIN matchmaking_ticket_user mtu ON mt.id = mtu.matchmaking_ticket_id
        JOIN matchmaking_user mu ON mtu.matchmaking_user_id = mu.id
        LEFT JOIN matchmaking_user_elo mue ON mu.id = mue.matchmaking_user_id
    WHERE mt.id = sqlc.arg(ticket_id)
    GROUP BY mue.matchmaking_arena_id
),
MatchElo AS (
    SELECT mm.id AS match_id,
        AVG(mue.elo) AS avg_elo
    FROM matchmaking_match mm
        JOIN matchmaking_ticket mt ON mm.id = mt.matchmaking_match_id
        JOIN matchmaking_ticket_user mtu ON mt.id = mtu.matchmaking_ticket_id
        JOIN matchmaking_user mu ON mtu.matchmaking_user_id = mu.id
        JOIN TicketElo te ON mm.matchmaking_arena_id = te.matchmaking_arena_id
        JOIN matchmaking_arena ma ON mm.matchmaking_arena_id = ma.id
        LEFT JOIN matchmaking_user_elo mue ON mu.id = mue.matchmaking_user_id
    WHERE mm.locked_at < NOW()
        AND te.player_count <= ma.max_players - COUNT(mu.id)
    GROUP BY mm.id
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
    JOIN MatchElo me ON mm.id = me.match_id
    JOIN TicketElo te ON me.avg_elo BETWEEN te.min_elo AND te.max_elo
    AND mm.matchmaking_arena_id = te.matchmaking_arena_id
ORDER BY ABS(me.avg_elo - AVG(mue.elo))
LIMIT 1;