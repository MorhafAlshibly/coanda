CREATE TABLE event (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    data JSON NOT NULL,
    started_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY (name)
) ENGINE = InnoDB;
CREATE TABLE event_round (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    event_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    scoring JSON NOT NULL,
    data JSON NOT NULL,
    ended_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE INDEX event_round_name_event_id_idx (name, event_id),
    UNIQUE INDEX event_round_ended_at_event_id_idx (ended_at, event_id),
    INDEX idx_event_id (event_id),
    CONSTRAINT fk_event_round_event FOREIGN KEY (event_id) REFERENCES event (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;
CREATE TABLE event_user (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    event_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    data JSON NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE INDEX event_user_user_id_event_id_idx (user_id, event_id),
    INDEX idx_event_id (event_id),
    CONSTRAINT fk_event_user_event FOREIGN KEY (event_id) REFERENCES event (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;
CREATE TABLE event_round_user (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    event_user_id BIGINT UNSIGNED NOT NULL,
    event_round_id BIGINT UNSIGNED NOT NULL,
    result BIGINT UNSIGNED NOT NULL,
    data JSON NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE INDEX event_round_event_user_id_event_round_id_idx (event_user_id, event_round_id),
    INDEX idx_event_user_id (event_user_id),
    INDEX idx_event_round_id (event_round_id),
    CONSTRAINT fk_event_round_user_event_round FOREIGN KEY (event_round_id) REFERENCES event_round (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_event_round_user_event_user FOREIGN KEY (event_user_id) REFERENCES event_user (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;
CREATE VIEW event_round_leaderboard AS WITH ranked_event_user AS (
    SELECT eru.id,
        eru.event_user_id,
        eru.event_round_id,
        eru.result,
        DENSE_RANK() OVER (
            PARTITION BY eru.event_round_id
            ORDER BY eru.result ASC
        ) AS ranking,
        eru.created_at,
        eru.updated_at
    FROM event_round_user eru
)
SELECT reu.id,
    reu.event_user_id,
    reu.event_round_id,
    reu.result,
    IF(
        reu.ranking <= JSON_LENGTH(er.scoring->'$.scoring'),
        JSON_UNQUOTE(
            JSON_EXTRACT(
                er.scoring,
                CONCAT('$.scoring[', reu.ranking - 1, ']')
            )
        ),
        '0'
    ) AS score,
    reu.ranking,
    reu.created_at,
    reu.updated_at
FROM ranked_event_user reu
    JOIN event_round er ON reu.event_round_id = er.id;
CREATE VIEW event_leaderboard AS WITH ranked_event_user AS (
    SELECT eru.id,
        eru.event_user_id,
        eru.event_round_id,
        eru.result,
        DENSE_RANK() OVER (
            PARTITION BY eru.event_round_id
            ORDER BY eru.result ASC
        ) AS ranking,
        eru.created_at,
        eru.updated_at,
        IF(
            DENSE_RANK() OVER (
                PARTITION BY eru.event_round_id
                ORDER BY eru.result ASC
            ) <= JSON_LENGTH(er.scoring->'$.scoring'),
            JSON_UNQUOTE(
                JSON_EXTRACT(
                    er.scoring,
                    CONCAT(
                        '$.scoring[',
                        DENSE_RANK() OVER (
                            PARTITION BY eru.event_round_id
                            ORDER BY eru.result ASC
                        ) - 1,
                        ']'
                    )
                )
            ),
            '0'
        ) AS score
    FROM event_round_user eru
        JOIN event_round er ON eru.event_round_id = er.id
),
user_scores AS (
    SELECT reu.event_user_id,
        SUM(CAST(reu.score AS UNSIGNED)) AS score
    FROM ranked_event_user reu
    GROUP BY reu.event_user_id
)
SELECT eu.id,
    eu.event_id,
    eu.user_id,
    us.score,
    DENSE_RANK() OVER (
        PARTITION BY eu.event_id
        ORDER BY us.score DESC
    ) AS ranking,
    eu.data,
    eu.created_at,
    eu.updated_at
FROM event_user eu
    JOIN user_scores us ON eu.id = us.event_user_id
ORDER BY us.score DESC;
CREATE VIEW event_with_round AS
SELECT e.id AS id,
    e.name AS name,
    current_round.id AS current_round_id,
    current_round.name AS current_round_name,
    e.data AS data,
    er.id AS round_id,
    er.name AS round_name,
    er.scoring AS round_scoring,
    er.data AS round_data,
    er.ended_at AS round_ended_at,
    er.created_at AS round_created_at,
    er.updated_at AS round_updated_at,
    e.started_at AS started_at,
    e.created_at AS created_at,
    e.updated_at AS updated_at
FROM event e
    LEFT JOIN event_round er ON e.id = er.event_id
    LEFT JOIN (
        SELECT er.id,
            er.name,
            er.event_id
        FROM event_round er
        WHERE er.ended_at > NOW()
        ORDER BY er.ended_at ASC
        LIMIT 1
    ) AS current_round ON e.id = current_round.event_id
ORDER BY e.id,
    er.id;