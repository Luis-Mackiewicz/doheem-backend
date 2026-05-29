-- name: GetGroupByID :one
SELECT * FROM groups
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListGroupsByUserID :many
SELECT g.* FROM groups g
JOIN group_members gm ON gm.group_id = g.id
WHERE gm.user_id = $1 AND g.deleted_at IS NULL AND gm.is_active = true
ORDER BY g.name;

-- name: CreateGroup :one
INSERT INTO groups (name, currency)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateGroup :one
UPDATE groups
SET name = COALESCE($2, name),
    currency = COALESCE($3, currency),
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteGroup :exec
UPDATE groups
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1;

-- name: DeactivateGroup :exec
UPDATE groups
SET is_active = false,
    inactive_since = NOW(),
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: ActivateGroup :exec
UPDATE groups
SET is_active = true,
    inactive_since = NULL,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;
