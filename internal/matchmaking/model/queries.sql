-- name: CreateArena :execresult
INSERT INTO matchmaking_arena (
        name,
        min_players,
        max_players,
        data
    )
VALUES (?, ?, ?, ?);