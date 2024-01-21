-- name: GetTeam :one
SELECT name,
  owner,
  score,
  ranking,
  data,
  created_at,
  updated_at
FROM ranked_team
WHERE name = sqlc.narg(name)
  OR owner = sqlc.narg(owner)
LIMIT 1;
-- name: GetTeamByMember :one
SELECT t.name,
  t.owner,
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
  score,
  ranking,
  data,
  created_at,
  updated_at
FROM ranked_team
ORDER BY score DESC
LIMIT ? OFFSET ?;
-- name: GetTeamMembers :many
SELECT team,
  user_id,
  data,
  joined_at,
  updated_at
FROM team_member
WHERE team = ?
ORDER BY joined_at ASC
LIMIT ? OFFSET ?;
-- name: GetTeamMember :one
SELECT team,
  user_id,
  data,
  joined_at,
  updated_at
FROM team_member
WHERE user_id = ?
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
-- name: DeleteTeam :exec
DELETE FROM team
WHERE name = sqlc.narg(name)
  OR owner = sqlc.narg(owner)
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
-- name: DeleteTeamMembersByTeam :exec
DELETE FROM team_member
WHERE team = ?;
-- name: DeleteTeamMembersByOwner :exec
DELETE FROM team_member
WHERE team = (
    SELECT name
    FROM team
    WHERE owner = ?
    LIMIT 1
  );
-- name: DeleteTeamMembersByMember :exec
DELETE FROM team_member
WHERE team = (
    SELECT team
    FROM team_member tm
    WHERE tm.user_id = ?
    LIMIT 1
  );
-- name: DeleteTeamOwner :exec
DELETE FROM team_owner
WHERE team = ?
  OR user_id = ?
LIMIT 1;
-- name: DeleteTeamOwnerByMember :exec
DELETE FROM team_owner
WHERE team = (
    SELECT team
    FROM team_member tm
    WHERE tm.user_id = ?
    LIMIT 1
  )
LIMIT 1;
-- name: UpdateTeamScore :exec
UPDATE team
SET score = score + ?
WHERE name = ?
  or owner = ?
LIMIT 1;
-- name: UpdateTeamData :exec
UPDATE team
SET data = ?
WHERE name = ?
  or owner = ?
LIMIT 1;
-- name: UpdateTeamScoreByMember :exec
UPDATE team
SET score = score + ?
WHERE name = (
    SELECT team
    FROM team_member
    WHERE user_id = ?
    LIMIT 1
  )
LIMIT 1;
-- name: UpdateTeamDataByMember :exec
UPDATE team
SET data = ?
WHERE name = (
    SELECT team
    FROM team_member
    WHERE user_id = ?
    LIMIT 1
  )
LIMIT 1;
-- name: UpdateTeamMemberData :exec
UPDATE team_member
SET data = ?
WHERE user_id = ?
LIMIT 1;
-- name: UpdateTeamMemberDataByTeam :exec
UPDATE team_member
SET data = ?
WHERE team = ?;
-- name: UpdateTeamMemberDataByOwner :exec
UPDATE team_member
SET data = ?
WHERE team = (
    SELECT name
    FROM team
    WHERE owner = ?
    LIMIT 1
  );
-- name: UpdateTeamMemberDataByMember :exec
UPDATE team_member
SET data = ?
WHERE team = (
    SELECT team
    FROM team_member tm
    WHERE tm.user_id = ?
    LIMIT 1
  );