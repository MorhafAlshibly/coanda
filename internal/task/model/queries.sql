-- name: CreateTask :execresult
INSERT INTO task (
        id,
        type,
        data,
        expires_at
    )
VALUES (?, ?, ?, ?);
-- name: GetTask :one
SELECT id,
    type,
    data,
    expires_at,
    completed_at,
    created_at,
    updated_at
FROM task
WHERE id = ?
    AND type = ?
    AND (
        expires_at IS NULL
        OR expires_at > NOW()
    )
LIMIT 1;
-- name: DeleteTask :execresult
DELETE FROM task
WHERE id = ?
    AND type = ?
    AND (
        expires_at IS NULL
        OR expires_at > NOW()
    )
LIMIT 1;
-- name: UpdateTask :execresult
UPDATE task
SET data = ?
WHERE id = ?
    AND type = ?
    AND (
        expires_at IS NULL
        OR expires_at > NOW()
    )
LIMIT 1;
-- name: CompleteTask :execresult
UPDATE task
SET completed_at = CURRENT_TIMESTAMP
WHERE id = ?
    AND type = ?
    AND (
        expires_at IS NULL
        OR expires_at > NOW()
    )
    AND completed_at IS NULL
LIMIT 1;