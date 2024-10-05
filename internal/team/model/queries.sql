-- name: GetTeams :many
SELECT id,
  name,
  score,
  ranking,
  data,
  created_at,
  updated_at,
  member_id,
  user_id,
  member_number,
  member_data,
  joined_at,
  member_updated_at,
  member_number_without_gaps
FROM ranked_team_with_member
WHERE member_number_without_gaps < CAST(sqlc.arg(member_limit) AS UNSIGNED)
  AND member_number_without_gaps >= CAST(sqlc.arg(member_offset) AS UNSIGNED)
ORDER BY score DESC,
  id,
  member_number
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