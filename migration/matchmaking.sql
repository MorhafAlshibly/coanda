CREATE TABLE matchmaking_arena (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) UNIQUE NOT NULL,
    min_players TINYINT UNSIGNED NOT NULL,
    max_players_per_ticket TINYINT UNSIGNED NOT NULL,
    max_players TINYINT UNSIGNED NOT NULL,
    data JSON NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE = InnoDB;
CREATE TABLE matchmaking_match (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    matchmaking_arena_id BIGINT UNSIGNED NOT NULL,
    private_server_id VARCHAR(255) NULL,
    data JSON NOT NULL,
    locked_at DATETIME NULL,
    started_at DATETIME NULL,
    ended_at DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_matchmaking_match_matchmaking_arena FOREIGN KEY (matchmaking_arena_id) REFERENCES matchmaking_arena (id) ON DELETE NO ACTION ON UPDATE CASCADE
) ENGINE = InnoDB;
CREATE TABLE matchmaking_ticket (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    matchmaking_match_id BIGINT UNSIGNED NULL,
    data JSON NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_matchmaking_ticket_matchmaking_match FOREIGN KEY (matchmaking_match_id) REFERENCES matchmaking_match (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;
CREATE TABLE matchmaking_user (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    matchmaking_ticket_id BIGINT UNSIGNED NULL,
    client_user_id BIGINT UNSIGNED UNIQUE NOT NULL,
    elo BIGINT NOT NULL,
    data JSON NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_matchmaking_user_matchmaking_ticket FOREIGN KEY (matchmaking_ticket_id) REFERENCES matchmaking_ticket (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;
CREATE TABLE matchmaking_ticket_arena (
    matchmaking_ticket_id BIGINT UNSIGNED NOT NULL,
    matchmaking_arena_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (matchmaking_ticket_id, matchmaking_arena_id),
    CONSTRAINT fk_matchmaking_ticket_arena_matchmaking_arena FOREIGN KEY (matchmaking_arena_id) REFERENCES matchmaking_arena (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_matchmaking_ticket_arena_matchmaking_ticket FOREIGN KEY (matchmaking_ticket_id) REFERENCES matchmaking_ticket (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;
CREATE VIEW matchmaking_ticket_with_user AS
SELECT mt.id AS ticket_id,
    mt.matchmaking_match_id,
    CASE
        WHEN mt.matchmaking_match_id IS NULL THEN "PENDING"
        WHEN mt.matchmaking_match_id IS NOT NULL
        AND (
            mm.ended_at > NOW()
            OR mm.ended_at IS NULL
        ) THEN "MATCHED"
        ELSE "ENDED"
    END AS status,
    COUNT(1) OVER (PARTITION BY mt.id) AS user_count,
    mt.data AS ticket_data,
    mt.created_at AS ticket_created_at,
    mt.updated_at AS ticket_updated_at,
    mu.id AS matchmaking_user_id,
    mu.client_user_id,
    mu.elo,
    DENSE_RANK() OVER (
        PARTITION BY mt.id
        ORDER BY mu.id
    ) AS user_number,
    mu.data AS user_data,
    mu.created_at AS user_created_at,
    mu.updated_at AS user_updated_at
FROM matchmaking_ticket mt
    JOIN matchmaking_user mu ON mu.matchmaking_ticket_id = mt.id
    LEFT JOIN matchmaking_match mm ON mt.matchmaking_match_id = mm.id
ORDER BY mt.id,
    mu.id;
CREATE VIEW matchmaking_ticket_with_user_and_arena AS
SELECT mtwu.*,
    ma.id AS arena_id,
    ma.name AS arena_name,
    ma.min_players AS arena_min_players,
    ma.max_players_per_ticket AS arena_max_players_per_ticket,
    ma.max_players AS arena_max_players,
    DENSE_RANK() OVER (
        PARTITION BY mtwu.ticket_id
        ORDER BY ma.id
    ) AS arena_number,
    ma.data AS arena_data,
    ma.created_at AS arena_created_at,
    ma.updated_at AS arena_updated_at
FROM matchmaking_ticket_with_user mtwu
    JOIN matchmaking_ticket_arena mta ON mtwu.ticket_id = mta.matchmaking_ticket_id
    JOIN matchmaking_arena ma ON mta.matchmaking_arena_id = ma.id
ORDER BY mtwu.ticket_id,
    mtwu.matchmaking_user_id,
    ma.id;
CREATE VIEW matchmaking_match_with_arena AS
SELECT mm.id AS match_id,
    mm.private_server_id,
    CASE
        WHEN mm.started_at IS NULL
        OR mm.started_at > NOW() THEN "PENDING"
        WHEN mm.ended_at IS NULL
        OR mm.ended_at > NOW() THEN "STARTED"
        ELSE "ENDED"
    END AS match_status,
    mm.data AS match_data,
    mm.locked_at,
    mm.started_at,
    mm.ended_at,
    mm.created_at AS match_created_at,
    mm.updated_at AS match_updated_at,
    ma.id AS arena_id,
    ma.name AS arena_name,
    ma.min_players AS arena_min_players,
    ma.max_players_per_ticket AS arena_max_players_per_ticket,
    ma.max_players AS arena_max_players,
    ma.data AS arena_data,
    ma.created_at AS arena_created_at,
    ma.updated_at AS arena_updated_at
FROM matchmaking_match mm
    JOIN matchmaking_arena ma ON mm.matchmaking_arena_id = ma.id
ORDER BY mm.id;
CREATE VIEW matchmaking_match_with_arena_and_ticket AS
SELECT mmwa.match_id,
    mmwa.private_server_id,
    mmwa.match_status,
    DENSE_RANK() OVER (
        PARTITION BY mmwa.match_id
        ORDER BY mtwua.ticket_id
    ) + DENSE_RANK() OVER (
        PARTITION BY mmwa.match_id
        ORDER BY mtwua.ticket_id DESC
    ) - 1 AS ticket_count,
    DENSE_RANK() OVER (
        PARTITION BY mmwa.match_id
        ORDER BY mtwua.matchmaking_user_id
    ) + DENSE_RANK() OVER (
        PARTITION BY mmwa.match_id
        ORDER BY mtwua.matchmaking_user_id DESC
    ) - 1 AS user_count,
    mmwa.match_data,
    mmwa.locked_at,
    mmwa.started_at,
    mmwa.ended_at,
    mmwa.match_created_at,
    mmwa.match_updated_at,
    mmwa.arena_id,
    mmwa.arena_name,
    mmwa.arena_min_players,
    mmwa.arena_max_players_per_ticket,
    mmwa.arena_max_players,
    mmwa.arena_data,
    mmwa.arena_created_at,
    mmwa.arena_updated_at,
    mtwua.ticket_id,
    mtwua.matchmaking_user_id,
    mtwua.status AS ticket_status,
    mtwua.user_count AS ticket_user_count,
    DENSE_RANK() OVER (
        PARTITION BY mmwa.match_id
        ORDER BY mtwua.ticket_id
    ) AS ticket_number,
    mtwua.ticket_data,
    mtwua.ticket_created_at,
    mtwua.ticket_updated_at,
    mtwua.client_user_id,
    mtwua.elo,
    mtwua.user_number,
    mtwua.user_data,
    mtwua.user_created_at,
    mtwua.user_updated_at,
    mtwua.arena_id AS ticket_arena_id,
    mtwua.arena_name AS ticket_arena_name,
    mtwua.arena_min_players AS ticket_arena_min_players,
    mtwua.arena_max_players_per_ticket AS ticket_arena_max_players_per_ticket,
    mtwua.arena_max_players AS ticket_arena_max_players,
    mtwua.arena_number,
    mtwua.arena_data AS ticket_arena_data,
    mtwua.arena_created_at AS ticket_arena_created_at,
    mtwua.arena_updated_at AS ticket_arena_updated_at
FROM matchmaking_match_with_arena mmwa
    LEFT JOIN matchmaking_ticket_with_user_and_arena mtwua ON mmwa.match_id = mtwua.matchmaking_match_id
ORDER BY mmwa.match_id,
    mtwua.ticket_id,
    mtwua.matchmaking_user_id,
    mtwua.arena_id;