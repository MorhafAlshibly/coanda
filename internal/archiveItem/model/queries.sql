-- name: GetExpiredItems :many
SELECT id,
    type,
    data,
    created_at,
    updated_at,
    expires_at
FROM item
WHERE expires_at < NOW()
    AND expires_at IS NOT NULL
ORDER BY id ASC
LIMIT ? OFFSET ?;
-- name: DeleteExpiredItems :execresult
DELETE FROM item
WHERE expires_at < NOW()
    AND expires_at IS NOT NULL
    AND id >= sqlc.arg(min_id)
    AND id <= sqlc.arg(max_id);