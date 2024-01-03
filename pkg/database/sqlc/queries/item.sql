-- name: CreateItem :execresult
INSERT INTO item (
  id, data, expires_at
) VALUES (
    ?, ?, ?
);

-- name: GetItem :one
SELECT * FROM item
WHERE id = ? LIMIT 1;

-- name: DeleteItem :exec
DELETE FROM item
WHERE id = ? LIMIT 1;

-- name: GetItems :many
SELECT * FROM item
ORDER BY created_at ASC
LIMIT ?
OFFSET ?;

-- name: UpdateItem :exec
UPDATE item
SET data = ?
WHERE id = ? LIMIT 1;



