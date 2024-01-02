-- name: GetRecord :one
SELECT * FROM ranked_record
WHERE name = ? AND user_id = ?;

-- name: GetRecords :many
SELECT * FROM ranked_record
WHERE name = ?
ORDER BY record ASC
LIMIT ?
OFFSET ?;

-- name: CreateRecord :execresult
INSERT INTO record (
  name, user_id, record, data
) VALUES (
  ?, ?, ?, ?
);

-- name: DeleteRecord :exec
DELETE FROM record
WHERE name = ? AND user_id = ?;