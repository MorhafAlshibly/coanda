-- name: CreateArena :execresult
INSERT INTO matchmaking_arena (
        name,
        min_players,
        max_players_per_ticket,
        max_players,
        data
    )
VALUES (?, ?, ?, ?, ?);
-- name: GetArenas :many
SELECT id,
    name,
    min_players,
    max_players_per_ticket,
    max_players,
    data,
    created_at,
    updated_at
FROM matchmaking_arena
ORDER BY created_at DESC
LIMIT ? OFFSET ?;
-- name: CreateMatchmakingUser :execresult
INSERT INTO matchmaking_user (client_user_id, elo, data)
VALUES (?, ?, ?);
-- name: GetMatchmakingUsers :many
SELECT id,
    client_user_id,
    elo,
    data,
    created_at,
    updated_at
FROM matchmaking_user
ORDER BY client_user_id ASC
LIMIT ? OFFSET ?;
-- name: CreateMatchmakingTicket :execresult
INSERT INTO matchmaking_ticket (data, elo_window, expires_at)
VALUES (?, ?, ?);
-- name: CreateMatchmakingTicketUser :execresult
INSERT INTO matchmaking_ticket_user (matchmaking_ticket_id, matchmaking_user_id)
VALUES (?, ?);
-- name: CreateMatchmakingTicketArena :execresult
INSERT INTO matchmaking_ticket_arena (matchmaking_ticket_id, matchmaking_arena_id)
VALUES (?, ?);
-- name: UpdateMatchmakingUserByClientUserId :execresult
UPDATE matchmaking_user
SET elo = ?,
    data = ?
WHERE client_user_id = ?;