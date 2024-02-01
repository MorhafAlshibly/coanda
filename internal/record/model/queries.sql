-- name: GetRecord :one
SELECT name,
  user_id,
  record,
  DENSE_RANK() OVER (
    PARTITION BY name
    ORDER BY record ASC
  ),
  data,
  created_at,
  updated_at
FROM record
WHERE name = ?
  AND user_id = ?
LIMIT 1;
-- name: GetRecords :many
SELECT name,
  user_id,
  record,
  DENSE_RANK() OVER (
    PARTITION BY name
    ORDER BY record ASC
  ),
  data,
  created_at,
  updated_at
FROM record
WHERE name = ?
  OR user_id = ?
ORDER BY record ASC
LIMIT ? OFFSET ?;
-- name: CreateRecord :execresult
INSERT INTO record (name, user_id, record, data)
VALUES (?, ?, ?, ?);
-- name: DeleteRecord :execresult
DELETE FROM record
WHERE name = ?
  AND user_id = ?
LIMIT 1;
-- name: UpdateRecordRecord :execresult
UPDATE record
SET record = ?
WHERE name = ?
  AND user_id = ?
LIMIT 1;
-- name: UpdateRecordData :execresult
UPDATE record
SET data = ?
WHERE name = ?
  AND user_id = ?
LIMIT 1;