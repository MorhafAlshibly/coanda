CREATE TABLE team (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    score BIGINT NOT NULL DEFAULT 0,
    data JSON NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE INDEX team_name_idx (name),
    INDEX team_score_idx (score DESC)
) ENGINE = InnoDB;
CREATE VIEW ranked_team AS
SELECT id,
    name,
    score,
    DENSE_RANK() OVER (
        ORDER BY score DESC
    ) AS ranking,
    data,
    created_at,
    updated_at
FROM team
ORDER BY score DESC;
CREATE TABLE team_member (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    team_id BIGINT UNSIGNED NOT NULL,
    member_number INT UNSIGNED NOT NULL,
    data JSON NOT NULL,
    joined_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE INDEX team_member_user_id_idx (user_id),
    UNIQUE INDEX team_member_team_id_member_number_idx (team_id, member_number),
    CONSTRAINT fk_team_member_team_id_is_team_id FOREIGN KEY (team_id) REFERENCES team(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;
CREATE VIEW ranked_team_with_member AS
SELECT t.id,
    t.name,
    t.score,
    t.ranking,
    t.data,
    t.created_at,
    t.updated_at,
    tm.id AS member_id,
    tm.user_id,
    tm.member_number,
    tm.data AS member_data,
    tm.joined_at,
    tm.updated_at AS member_updated_at,
    ROW_NUMBER() OVER (
        PARTITION BY t.id
        ORDER BY tm.member_number
    ) AS member_number_without_gaps
FROM ranked_team t
    LEFT JOIN team_member tm ON t.id = tm.team_id
ORDER BY t.score DESC,
    t.id,
    tm.member_number;
CREATE VIEW team_with_first_open_member AS
SELECT t.id,
    COALESCE(
        (
            -- Find the first gap in the sequence for the team
            SELECT MIN(m1.member_number + 1)
            FROM team_member m1
            WHERE m1.team_id = t.id
                AND NOT EXISTS (
                    SELECT 1
                    FROM team_member m2
                    WHERE m2.team_id = t.id
                        AND m2.member_number = m1.member_number + 1
                )
        ),
        -- If no members exist, default to 1
        (
            SELECT 1
            WHERE NOT EXISTS (
                    SELECT 1
                    FROM team_member m
                    WHERE m.team_id = t.id
                )
        )
    ) AS first_open_member
FROM team t;