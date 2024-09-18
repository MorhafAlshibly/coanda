-- name: GetTeams :many
SELECT name,
  owner,
  score,
  ranking,
  data,
  created_at,
  updated_at
FROM ranked_team
ORDER BY score DESC
LIMIT ? OFFSET ?;
-- name: SearchTeams :many
SELECT name,
  owner,
  score,
  ranking,
  data,
  created_at,
  updated_at
FROM ranked_team
WHERE name LIKE CONCAT('%', sqlc.arg(query), '%')
ORDER BY score DESC
LIMIT ? OFFSET ?;
-- name: GetTeamMember :one
SELECT team,
  user_id,
  member_number,
  data,
  joined_at,
  updated_at
FROM team_member
WHERE user_id = sqlc.arg(member)
LIMIT 1;
-- name: GetHighestMemberNumber :one
SELECT max_member_number
FROM last_team_member
WHERE team = sqlc.arg(team)
LIMIT 1;
-- name: CreateTeam :execresult
INSERT INTO team (name, owner, score, data)
VALUES (?, ?, ?, ?);
-- name: CreateTeamOwner :execresult
INSERT INTO team_owner (team, user_id)
VALUES (?, ?);
-- name: CreateTeamMember :execresult
INSERT INTO team_member (team, user_id, member_number, data)
VALUES (?, ?, ?, ?);
-- name: DeleteTeamMember :execresult
DELETE FROM team_member
WHERE user_id = ?
LIMIT 1;
-- name: UpdateTeamMember :execresult
UPDATE team_member
SET data = ?
WHERE user_id = ?
LIMIT 1;