-- name: GetTeams :many
SELECT id,
  name,
  score,
  ranking,
  data,
  created_at,
  updated_at
FROM ranked_team
ORDER BY score DESC
LIMIT ? OFFSET ?;
-- name: SearchTeams :many
SELECT id,
  name,
  score,
  ranking,
  data,
  created_at,
  updated_at
FROM ranked_team
WHERE name LIKE CONCAT('%', sqlc.arg(query), '%')
ORDER BY score DESC
LIMIT ? OFFSET ?;
-- name: GetFirstOpenMemberNumber :one
SELECT first_open_member
FROM team_with_first_open_member
WHERE id = sqlc.arg(team)
LIMIT 1;
-- name: CreateTeam :execresult
INSERT INTO team (name, score, data)
VALUES (?, ?, ?);
-- name: CreateTeamMember :execresult
INSERT INTO team_member (user_id, team_id, member_number, data)
VALUES (?, ?, ?, ?);