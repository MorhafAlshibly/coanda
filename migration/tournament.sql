CREATE TABLE tournament (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    tournament_interval ENUM('daily', 'weekly', 'monthly', 'unlimited') NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    score BIGINT NOT NULL DEFAULT 0,
    data JSON NOT NULL,
    tournament_started_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE INDEX name_tournament_interval_user_id_tournament_started_at_idx (
        name ASC,
        tournament_interval ASC,
        user_id ASC,
        tournament_started_at DESC
    ),
    INDEX name_tournament_interval_score_tournament_started_at_idx (
        name ASC,
        tournament_interval ASC,
        score DESC,
        tournament_started_at DESC
    ),
    INDEX user_id_score_tournament_started_at_idx (
        user_id ASC,
        score DESC,
        tournament_started_at DESC
    )
) ENGINE = InnoDB;
CREATE VIEW ranked_tournament AS
SELECT id,
    name,
    tournament_interval,
    user_id,
    score,
    DENSE_RANK() OVER (
        PARTITION BY name,
        tournament_interval,
        tournament_started_at
        ORDER BY score DESC
    ) AS ranking,
    data,
    tournament_started_at,
    created_at,
    updated_at
FROM tournament
ORDER BY name ASC,
    tournament_interval ASC,
    score DESC,
    tournament_started_at DESC;
CREATE TABLE archived_tournament (
    id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    tournament_interval ENUM('daily', 'weekly', 'monthly', 'unlimited') NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    score BIGINT NOT NULL DEFAULT 0,
    data JSON NOT NULL,
    tournament_started_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    archived_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE INDEX a_name_tournament_interval_user_id_tournament_started_at_idx (
        name ASC,
        tournament_interval ASC,
        user_id ASC,
        tournament_started_at DESC
    ),
    PRIMARY KEY (id)
) ENGINE = InnoDB;