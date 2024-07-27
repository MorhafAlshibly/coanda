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
  member_number,
  data,
  joined_at,
  updated_at
FROM team_member
WHERE user_id = sqlc.arg(member)
LIMIT 1;
-- name: GetHighestMemberNumber :one
SELECT MAX(member_number) AS member_number
FROM team_member
WHERE team = sqlc.arg(team);
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