CREATE TABLE matchmaking_user (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT UNSIGNED UNIQUE NOT NULL,
    data JSON NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE = InnoDB;
CREATE TABLE matchmaking_arena (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) UNIQUE NOT NULL,
    min_players TINYINT UNSIGNED NOT NULL,
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
    matchmaking_match_id INT NULL,
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
    mu.user_id,
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