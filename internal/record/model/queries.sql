-- name: GetRecord :one
SELECT *
FROM ranked_record
WHERE name = ?
  AND user_id = ?
LIMIT 1;
-- name: GetRecordsByName :many
SELECT *
FROM ranked_record
WHERE name = ?
ORDER BY record ASC
LIMIT ? OFFSET ?;
-- name: GetRecordsByUser :many
SELECT *
FROM ranked_record
WHERE user_id = ?
ORDER BY record ASC
LIMIT ? OFFSET ?;
-- name: GetRecords :many
SELECT *
FROM ranked_record
ORDER BY record ASC
LIMIT ? OFFSET ?;
-- name: CreateRecord :execresult
INSERT INTO record (name, user_id, record, data)
VALUES (?, ?, ?, ?);
-- name: DeleteRecord :exec
DELETE FROM record
WHERE name = ?
  AND user_id = ?
LIMIT 1;
-- name: UpdateRecordRecord :exec
UPDATE record
SET record = ?
WHERE name = ?
  AND user_id = ?
LIMIT 1;
-- name: UpdateRecordData :exec
UPDATE record
SET data = ?
WHERE name = ?
  AND user_id = ?
LIMIT 1;
-- name: UpdateRecord :exec
UPDATE record
SET record = ?,
  data = ?
WHERE name = ?
  AND user_id = ?
LIMIT 1;