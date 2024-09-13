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
    expires_at,
    created_at,
    updated_at
FROM item
WHERE id = ?
    AND type = ?
    AND (
        expires_at IS NULL
        OR expires_at > NOW()
    )
LIMIT 1;
-- name: DeleteItem :execresult
DELETE FROM item
WHERE id = ?
    AND type = ?
    AND (
        expires_at IS NULL
        OR expires_at > NOW()
    )
LIMIT 1;
-- name: UpdateItem :execresult
UPDATE item
SET data = ?
WHERE id = ?
    AND type = ?
    AND (
        expires_at IS NULL
        OR expires_at > NOW()
    )
LIMIT 1;