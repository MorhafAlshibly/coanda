CREATE TABLE tournament (
    name VARCHAR(255) NOT NULL,
    tournament_interval ENUM('daily', 'weekly', 'monthly', 'unlimited') NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    score BIGINT NOT NULL DEFAULT 0,
    data JSON NOT NULL,
    tournament_started_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (
        name,
        tournament_interval,
        user_id,
        tournament_started_at
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
SELECT name,
    tournament_interval,
    user_id,
    score,
    DENSE_RANK() OVER (
        PARTITION BY name,
        tournament_interval
        ORDER BY score DESC
    ) AS ranking,
    data,
    tournament_started_at,
    created_at,
    updated_at
FROM tournament
ORDER BY name ASC,
    tournament_interval ASC,
    score DESC;