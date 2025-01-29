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
    matchmaking_ticket_id,
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
-- name: AddTicketIDToUser :execresult
UPDATE matchmaking_user
SET matchmaking_ticket_id = ?
WHERE id = ?
    AND matchmaking_ticket_id IS NULL;
-- name: CreateMatchmakingTicketArena :execresult
INSERT INTO matchmaking_ticket_arena (matchmaking_ticket_id, matchmaking_arena_id)
VALUES (?, ?);
-- name: DeleteAllExpiredTickets :execresult
DELETE FROM matchmaking_ticket
WHERE expires_at < NOW()
    AND matchmaking_match_id IS NULL;