-- name: GetTeamByName :one
SELECT * FROM ranked_team
WHERE name = ? LIMIT 1;

-- name: GetTeamByOwner :one
SELECT * FROM ranked_team
WHERE owner = ? LIMIT 1;

-- name: GetTeamByMember :one
SELECT * FROM ranked_team_members
WHERE member = ? LIMIT 1;

-- name: GetTeams :many
SELECT * FROM ranked_team
ORDER BY score DESC
LIMIT ?
OFFSET ?;

-- name: CreateTeam :execresult
INSERT INTO team (
  name, owner, score, data
) VALUES (
  ?, ?, ?, ?
);

-- name: CreateTeamMember :exec
START TRANSACTION;
SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;
INSERT INTO team_members (team_name, user_id)
SELECT DISTINCT ?, ?
FROM team_members tm_outer
WHERE (SELECT COUNT(*) FROM team_members tm_inner WHERE tm_inner.team_name = tm_outer.team_name FOR UPDATE) < ?;
COMMIT;

-- name: DeleteTeam :exec
DELETE FROM team
WHERE name = ? LIMIT 1;

-- name: DeleteTeamMember :exec
DELETE FROM team_members
WHERE team_name = ? AND user_id = ? LIMIT 1;

-- name: UpdateTeam :exec
UPDATE team
SET score = ?, data = ?
WHERE name = ? LIMIT 1;