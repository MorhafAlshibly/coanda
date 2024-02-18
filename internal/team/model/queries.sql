-- name: GetTeams :many
SELECT *
FROM ranked_team
ORDER BY score DESC
LIMIT ? OFFSET ?;
-- name: SearchTeams :many
SELECT *
FROM ranked_team
WHERE name LIKE CONCAT('%', sqlc.arg(query), '%')
ORDER BY score DESC
LIMIT ? OFFSET ?;
-- name: GetTeamMember :one
SELECT team,
  user_id,
  data,
  joined_at,
  updated_at
FROM team_member
WHERE user_id = sqlc.arg(member)
LIMIT 1;
-- name: CreateTeam :execresult
INSERT INTO team (name, owner, score, data)
VALUES (?, ?, ?, ?);
-- name: CreateTeamOwner :execresult
INSERT INTO team_owner (team, user_id)
VALUES (?, ?);
-- name: CreateTeamMember :execresult
INSERT INTO team_member (team, user_id, data)
SELECT sqlc.arg(team),
  ?,
  ?
FROM dual
WHERE (
    SELECT COUNT(*)
    FROM team_member tm
    WHERE tm.team = sqlc.arg(team)
  ) < CAST(sqlc.arg(max_members) as unsigned);
-- name: DeleteTeamMember :execresult
DELETE FROM team_member
WHERE user_id = ?
LIMIT 1;
-- name: UpdateTeamMember :execresult
UPDATE team_member
SET data = ?
WHERE user_id = ?
LIMIT 1;