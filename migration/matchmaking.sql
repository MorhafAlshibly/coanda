CREATE TABLE matchmaking_user (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    client_user_id BIGINT UNSIGNED UNIQUE NOT NULL,
    data JSON NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE = InnoDB;
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
    matchmaking_arena_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    data JSON NOT NULL,
    locked_at DATETIME NULL,
    started_at DATETIME NULL,
    ended_at DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_matchmaking_match_matchmaking_arena FOREIGN KEY (matchmaking_arena_id) REFERENCES matchmaking_arena (id) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE = InnoDB;
CREATE TABLE matchmaking_ticket (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    matchmaking_match_id BIGINT UNSIGNED NULL,
    data JSON NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_matchmaking_ticket_matchmaking_match FOREIGN KEY (matchmaking_match_id) REFERENCES matchmaking_match (id) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE = InnoDB;
CREATE TABLE matchmaking_ticket_arena (
    matchmaking_ticket_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    matchmaking_arena_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (matchmaking_ticket_id, matchmaking_arena_id),
    CONSTRAINT fk_matchmaking_ticket_arena_matchmaking_arena FOREIGN KEY (matchmaking_arena_id) REFERENCES matchmaking_arena (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT fk_matchmaking_ticket_arena_matchmaking_ticket FOREIGN KEY (matchmaking_ticket_id) REFERENCES matchmaking_ticket (id) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE = InnoDB;
CREATE TABLE matchmaking_user_elo (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    elo INT NOT NULL,
    matchmaking_user_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    matchmaking_arena_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (id),
    CONSTRAINT fk_matchmaking_user_elo_matchmaking_user FOREIGN KEY (matchmaking_user_id) REFERENCES matchmaking_user (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT fk_matchmaking_user_elo_matchmaking_arena FOREIGN KEY (matchmaking_arena_id) REFERENCES matchmaking_arena (id) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE = InnoDB;
CREATE TABLE matchmaking_ticket_user (
    matchmaking_ticket_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    matchmaking_user_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (matchmaking_ticket_id, matchmaking_user_id),
    CONSTRAINT fk_matchmaking_ticket_user_matchmaking_ticket FOREIGN KEY (matchmaking_ticket_id) REFERENCES matchmaking_ticket (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT fk_matchmaking_ticket_user_matchmaking_user FOREIGN KEY (matchmaking_user_id) REFERENCES matchmaking_user (id) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE = InnoDB;
CREATE VIEW matchmaking_user_with_elo AS
SELECT mu.id,
    mu.client_user_id,
    JSON_ARRAYAGG(
        JSON_OBJECT(
            'arena_id',
            mue.matchmaking_arena_id,
            'elo',
            mue.elo
        )
    ) AS elos,
    mu.data,
    mu.created_at,
    mu.updated_at
FROM matchmaking_user mu
    LEFT JOIN matchmaking_user_elo mue ON mu.id = mue.matchmaking_user_id
GROUP BY mu.id;
CREATE VIEW matchmaking_ticket_with_user AS
SELECT mt.id,
    mu.id AS matchmaking_user_id,
    mu.client_user_id,
    mu.data AS user_data,
    mu.created_at AS user_created_at,
    mu.updated_at AS user_updated_at,
    mt.matchmaking_match_id,
    CASE
        WHEN mt.matchmaking_match_id IS NULL
        AND mt.expires_at > NOW() THEN "PENDING"
        WHEN mt.matchmaking_match_id IS NULL
        AND mt.expires_at < NOW() THEN "EXPIRED"
        WHEN mt.matchmaking_match_id IS NOT NULL
        AND mm.ended_at > NOW() THEN "MATCHED"
        ELSE "ENDED"
    END AS status,
    mt.data AS ticket_data,
    mt.expires_at,
    mt.created_at AS ticket_created_at,
    mt.updated_at AS ticket_updated_at
FROM matchmaking_ticket mt
    JOIN matchmaking_ticket_user mtu ON mt.id = mtu.matchmaking_ticket_id
    JOIN matchmaking_user mu ON mtu.matchmaking_user_id = mu.id
    LEFT JOIN matchmaking_match mm ON mt.matchmaking_match_id = mm.id
GROUP BY mt.id;
CREATE VIEW matchmaking_ticket_with_user_and_arena AS
SELECT mt.id,
    mu.id AS matchmaking_user_id,
    mu.client_user_id,
    JSON_ARRAYAGG(
        JSON_OBJECT(
            'arena_id',
            mue.matchmaking_arena_id,
            'elo',
            mue.elo
        )
    ) AS elos,
    mu.data AS user_data,
    mu.created_at AS user_created_at,
    mu.updated_at AS user_updated_at,
    JSON_ARRAYAGG(
        JSON_OBJECT(
            'arena_id',
            mta.matchmaking_arena_id,
            'name',
            ma.name,
            'min_players',
            ma.min_players,
            'max_players_per_ticket',
            ma.max_players_per_ticket,
            'max_players',
            ma.max_players,
            'data',
            ma.data,
            'created_at',
            ma.created_at,
            'updated_at',
            ma.updated_at
        )
    ) AS arenas,
    mt.matchmaking_match_id,
    CASE
        WHEN mt.matchmaking_match_id IS NULL
        AND mt.expires_at > NOW() THEN "PENDING"
        WHEN mt.matchmaking_match_id IS NULL
        AND mt.expires_at < NOW() THEN "EXPIRED"
        WHEN mt.matchmaking_match_id IS NOT NULL
        AND mm.ended_at > NOW() THEN "MATCHED"
        ELSE "ENDED"
    END AS status,
    mt.data AS ticket_data,
    mt.expires_at,
    mt.created_at AS ticket_created_at,
    mt.updated_at AS ticket_updated_at
FROM matchmaking_ticket mt
    JOIN matchmaking_ticket_user mtu ON mt.id = mtu.matchmaking_ticket_id
    JOIN matchmaking_user mu ON mtu.matchmaking_user_id = mu.id
    LEFT JOIN matchmaking_user_elo mue ON mu.id = mue.matchmaking_user_id
    LEFT JOIN matchmaking_ticket_arena mta ON mt.id = mta.matchmaking_ticket_id
    LEFT JOIN matchmaking_arena ma ON mta.matchmaking_arena_id = ma.id
    LEFT JOIN matchmaking_match mm ON mt.matchmaking_match_id = mm.id
GROUP BY mt.id;