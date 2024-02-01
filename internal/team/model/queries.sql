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
-- name: GetTeamMembersByOwner :many
SELECT tm.team,
  tm.user_id,
  tm.data,
  tm.joined_at,
  tm.updated_at
FROM team_member tm
  JOIN team t ON tm.team = t.name
WHERE t.owner = ?
ORDER BY tm.joined_at ASC
LIMIT ? OFFSET ?;
-- name: GetTeamMembersByMember :many
SELECT tm.team,
  tm.user_id,
  tm.data,
  tm.joined_at,
  tm.updated_at
FROM team_member tm
  JOIN team_member tm2 ON tm.team = tm2.team
WHERE tm2.user_id = ?
ORDER BY tm.joined_at ASC
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
-- name: DeleteTeam :execresult
DELETE FROM team
WHERE name = sqlc.narg(name)
  OR owner = sqlc.narg(owner)
LIMIT 1;
-- name: DeleteTeamByMember :execresult
DELETE FROM team
WHERE name = (
    SELECT team
    FROM team_member
    WHERE user_id = ?
    LIMIT 1
  )
LIMIT 1;
-- name: DeleteTeamMember :execresult
DELETE FROM team_member
WHERE user_id = ?
LIMIT 1;
-- name: DeleteTeamMembersByTeam :execresult
DELETE FROM team_member
WHERE team = ?;
-- name: DeleteTeamMembersByOwner :execresult
DELETE FROM team_member
WHERE team = (
    SELECT name
    FROM team
    WHERE owner = ?
    LIMIT 1
  );
-- name: DeleteTeamMembersByMember :execresult
DELETE FROM team_member
WHERE team = (
    SELECT team
    FROM team_member tm
    WHERE tm.user_id = ?
    LIMIT 1
  );
-- name: DeleteTeamOwner :execresult
DELETE FROM team_owner
WHERE team = ?
  OR user_id = ?
LIMIT 1;
-- name: DeleteTeamOwnerByMember :execresult
DELETE FROM team_owner
WHERE team = (
    SELECT team
    FROM team_member tm
    WHERE tm.user_id = ?
    LIMIT 1
  )
LIMIT 1;
-- name: UpdateTeam :execresult
UPDATE team
SET score = CASE
    WHEN sqlc.narg(score) IS NOT NULL THEN sqlc.narg(score) + CASE
      WHEN sqlc.arg(increment_score) IS NOT NULL THEN score
      ELSE 0
    END
    ELSE score
  END,
  data = CASE
    WHEN CAST(sqlc.arg(data_exists) as unsigned) != 0 THEN sqlc.arg(data)
    ELSE data
  END
WHERE name = sqlc.narg(name)
  or name = CASE
    WHEN CAST(sqlc.narg(member) as unsigned) IS NOT NULL THEN (
      SELECT team
      FROM team_member
      WHERE user_id = sqlc.narg(member)
      LIMIT 1
    )
    ELSE NULL
  END
  or owner = sqlc.narg(owner)
LIMIT 1;
-- name: UpdateTeamMember :execresult
UPDATE team_member
SET data = ?
WHERE user_id = ?;
LIMIT 1;
-- name: UpdateTeamMembers :execresult
UPDATE team_member tm
SET data = ?
WHERE tm.team = sqlc.narg(name)
  OR tm.team = CASE
    WHEN CAST(sqlc.narg(owner) as unsigned) IS NOT NULL THEN (
      SELECT name
      FROM team
      WHERE owner = sqlc.narg(owner)
      LIMIT 1
    )
    ELSE NULL
  END
  OR tm.team = CASE
    WHEN CAST(sqlc.narg(member) as unsigned) IS NOT NULL THEN (
      SELECT team
      FROM team_member tm_user_id
      WHERE tm_user_id.user_id = sqlc.narg(member)
      LIMIT 1
    )
    ELSE NULL
  END;