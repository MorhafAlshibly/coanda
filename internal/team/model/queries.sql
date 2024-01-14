-- name: GetTeamByNameOrOwner :one
SELECT t.name,
  t.owner,
  tm.user_id AS members,
  t.score,
  DENSE_RANK() OVER (
    ORDER BY t.score DESC
  ) AS ranking,
  t.data,
  t.created_at,
  t.updated_at
FROM team t
  JOIN team_member tm ON t.name = tm.team
WHERE t.name = ?
  OR t.owner = ?
GROUP BY tm.user_id
ORDER BY score DESC;
-- name: GetTeamByMember :one
SELECT t.name,
  t.owner,
  JSON_ARRAYAGG(t.members) as members,
  t.score,
  t.ranking,
  t.data,
  t.created_at,
  t.updated_at
FROM ranked_team t
  JOIN team_member tm ON t.name = tm.team
WHERE tm.user_id = ?
LIMIT 1;
-- name: GetTeams :many
SELECT name,
  owner,
  JSON_ARRAYAGG(members) as members,
  score,
  ranking,
  data,
  created_at,
  updated_at
FROM ranked_team
GROUP BY name
ORDER BY score DESC
LIMIT ? OFFSET ?;
-- name: CreateTeam :execresult
INSERT INTO team (name, owner, score, data)
VALUES (?, ?, ?, ?);
-- name: CreateTeamOwner :execresult
INSERT INTO team_owner (team, user_id)
VALUES (?, ?);
-- name: CreateTeamMember :execresult
START TRANSACTION;
SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;
INSERT INTO team_member (team, user_id)
SELECT DISTINCT ?,
  ?
FROM team_member tm_outer
WHERE (
    SELECT COUNT(*)
    FROM team_member tm_inner
    WHERE tm_inner.team = tm_outer.team FOR
    UPDATE
  ) < ?;
COMMIT;
-- name: DeleteTeamByNameOrOwner :exec
DELETE FROM team
WHERE name = ?
  OR owner = ?
LIMIT 1;
-- name: DeleteTeamByMember :exec
DELETE FROM team
WHERE name = (
    SELECT team
    FROM team_member
    WHERE user_id = ?
    LIMIT 1
  )
LIMIT 1;
-- name: DeleteTeamMember :exec
DELETE FROM team_member
WHERE user_id = ?
LIMIT 1;
-- name: DeleteTeamOwner :exec
DELETE FROM team_owner
WHERE team = ?
  OR user_id = ?
LIMIT 1;
-- name: UpdateTeamByNameOrOwner :exec
UPDATE team
SET score = score + CASE
    WHEN ? THEN ?
    ELSE 0
  END,
  data = CASE
    WHEN ? THEN ?
    ELSE data
  END
WHERE name = ?
  OR owner = ?
LIMIT 1;
-- name: UpdateTeamByMember :exec
UPDATE team
SET score = score + CASE
    WHEN ? THEN ?
    ELSE 0
  END,
  data = CASE
    WHEN ? THEN ?
    ELSE data
  END
WHERE name = (
    SELECT team
    FROM team_member
    WHERE user_id = ?
    LIMIT 1
  )
LIMIT 1;