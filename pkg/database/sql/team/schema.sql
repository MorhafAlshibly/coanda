CREATE TABLE team (
    name VARCHAR(255) NOT NULL,
    owner BIGINT UNSIGNED NOT NULL,
    score BIGINT NOT NULL,
    data JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (name),
    UNIQUE (owner)
);

CREATE INDEX idx_owner ON team (owner);
CREATE INDEX idx_score ON team (score);

CREATE TABLE team_members (
    team_name VARCHAR(255) NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (user_id),
    FOREIGN KEY (team_name) REFERENCES team(name),
    UNIQUE (user_id)
);

CREATE INDEX idx_team_name ON team_members (team_name);

CREATE VIEW ranked_team_members AS
SELECT
    t.name AS name,
    t.owner AS owner,
    tm.user_id AS member,
    t.score AS score,
    t.data AS data,
    DENSE_RANK() OVER (ORDER BY t.score DESC) AS ranking,
    t.created_at AS created_at,
    t.updated_at AS updated_at
FROM
    team t
JOIN
    team_members tm ON t.name = tm.team_name;

CREATE VIEW ranked_team AS
SELECT
    name,
    owner,
    GROUP_CONCAT(members) AS members,
    score,
    data,
    ranking,
    created_at,
    updated_at
FROM
    ranked_team_members;