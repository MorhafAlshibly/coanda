CREATE TABLE team (
    name VARCHAR(255) PRIMARY KEY NOT NULL,
    owner BIGINT UNSIGNED UNIQUE NOT NULL,
    score BIGINT NOT NULL DEFAULT 0,
    data JSON NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX team_score_idx (score DESC)
) ENGINE = InnoDB;
CREATE VIEW ranked_team AS
SELECT name,
    owner,
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
    team VARCHAR(255) NOT NULL,
    user_id BIGINT UNSIGNED PRIMARY KEY NOT NULL,
    member_number INT UNSIGNED NOT NULL,
    data JSON NOT NULL,
    joined_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX team_member_number_idx (team ASC, member_number ASC),
    CONSTRAINT fk_team_member_team_is_team_name FOREIGN KEY (team) REFERENCES team(name) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;
CREATE TABLE team_owner (
    team VARCHAR(255) PRIMARY KEY NOT NULL,
    user_id BIGINT UNSIGNED UNIQUE NOT NULL,
    CONSTRAINT fk_team_owner_team_is_team_name FOREIGN KEY (team) REFERENCES team(name) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_team_owner_user_id_is_team_member_user_id FOREIGN KEY (user_id) REFERENCES team_member(user_id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB;
CREATE VIEW last_team_member AS
SELECT team,
    MAX(member_number) AS max_member_number
FROM team_member
GROUP BY team;