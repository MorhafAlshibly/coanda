-- name: CreateArena :execresult
INSERT INTO matchmaking_arena (
        name,
        min_players,
        max_players,
        data
    )
VALUES (?, ?, ?, ?);
-- name: GetArenas :many
SELECT id,
    name,
    min_players,
    max_players,
    data,
    created_at,
    updated_at
FROM matchmaking_arena
ORDER BY created_at DESC
LIMIT ? OFFSET ?;
-- name: CreateMatchmakingUser :execresult
INSERT INTO matchmaking_user (user_id, data)
VALUES (?, ?);
-- name: GetMatchmakingUsers :many
SELECT id,
    user_id,
    elos,
    data,
    created_at,
    updated_at
FROM matchmaking_user_with_elo
ORDER BY user_id ASC
LIMIT ? OFFSET ?;