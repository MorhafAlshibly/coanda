-- name: CreateItem :execresult
INSERT INTO item (
        id,
        type,
        data,
        expires_at
    )
VALUES (?, ?, ?, ?);
-- name: GetItem :one
SELECT id,
    type,
    data,
    created_at,
    updated_at,
    expires_at
FROM item
WHERE id = ?
    AND type = ?
LIMIT 1;
-- name: GetItems :many
SELECT id,
    type,
    data,
    created_at,
    updated_at,
    expires_at
FROM item
WHERE type = CASE
        WHEN sqlc.narg(type) IS NOT NULL THEN sqlc.narg(type)
        ELSE type
    END
ORDER BY id ASC
LIMIT ? OFFSET ?;
-- name: DeleteItem :execresult
DELETE FROM item
WHERE id = ?
    AND type = ?
LIMIT 1;
-- name: UpdateItem :execresult
UPDATE item
SET data = CASE
        WHEN CAST(sqlc.arg(data_exists) as unsigned) != 0 THEN sqlc.narg(data)
        ELSE data
    END,
    expires_at = CASE
        WHEN sqlc.narg(expires_at) IS NOT NULL THEN sqlc.narg(expires_at)
        ELSE expires_at
    END
WHERE id = ?
    AND type = ?
LIMIT 1;