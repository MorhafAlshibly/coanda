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
INSERT INTO matchmaking_user (client_user_id, data)
VALUES (?, ?);
-- name: GetMatchmakingUsers :many
SELECT id,
    client_user_id,
    elos,
    data,
    created_at,
    updated_at
FROM matchmaking_user_with_elo
ORDER BY client_user_id ASC
LIMIT ? OFFSET ?;
-- name: CreateMatchmakingTicket :execresult
INSERT INTO matchmaking_ticket (data, elo_window, expires_at)
SELECT ?,
    0,
    ?
FROM DUAL
WHERE NOT EXISTS (
        SELECT 1
        FROM matchmaking_ticket_user mtu
            JOIN matchmaking_ticket mt ON mtu.matchmaking_ticket_id = mt.id
            LEFT JOIN matchmaking_match mm ON mt.matchmaking_match_id = mm.id
        WHERE FIND_IN_SET(
                mtu.matchmaking_user_id,
                sqlc.arg(ids_separated_by_comma)
            )
            AND (
                (
                    mt.matchmaking_match_id IS NULL
                    AND mt.expires_at > NOW()
                )
                OR (
                    mt.matchmaking_match_id IS NOT NULL
                    AND mm.ended_at > NOW()
                )
            )
    );
-- name: CreateMatchmakingTicketUser :execresult
INSERT INTO matchmaking_ticket_user (matchmaking_ticket_id, matchmaking_user_id)
VALUES (?, ?);
-- name: CreateMatchmakingTicketArena :execresult
INSERT INTO matchmaking_ticket_arena (matchmaking_ticket_id, matchmaking_arena_id)
VALUES (?, ?);
-- name: CreateMatchmakingUserElo :execresult
INSERT INTO matchmaking_user_elo (matchmaking_user_id, matchmaking_arena_id, elo)
VALUES (?, ?, ?);
-- name: SetAllMatchmakingUserElos :execresult
INSERT INTO matchmaking_user_elo (matchmaking_user_id, matchmaking_arena_id, elo)
SELECT mu.id,
    ma.id,
    ?
FROM matchmaking_user mu
    JOIN matchmaking_arena ma
WHERE mu.id = ?;