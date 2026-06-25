-- name: GetGroupByID :one
SELECT id, name, description, monthly_fee, photo_url, invite_token, created_at, updated_at FROM groups
WHERE id = $1;

-- name: ListGroupsByUserID :many
SELECT g.id, g.name, g.description, g.monthly_fee, g.photo_url, g.invite_token, g.created_at, g.updated_at FROM groups g
JOIN group_members gm ON gm.group_id = g.id
WHERE gm.user_id = $1
ORDER BY g.name;

-- name: CreateGroup :one
INSERT INTO groups (name, description, invite_token)
VALUES ($1, $2, gen_random_uuid()::text)
RETURNING id, name, description, monthly_fee, photo_url, invite_token, created_at, updated_at;

-- name: UpdateGroup :one
UPDATE groups
SET name = COALESCE($2, name),
    description = COALESCE($3, description),
    monthly_fee = COALESCE($4, monthly_fee),
    photo_url = COALESCE($5, photo_url),
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, description, monthly_fee, photo_url, invite_token, created_at, updated_at;

-- name: RegenerateInviteToken :one
UPDATE groups
SET invite_token = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, description, monthly_fee, photo_url, invite_token, created_at, updated_at;

-- name: DeleteGroup :exec
DELETE FROM groups
WHERE id = $1;

-- name: CountGroupsByUserID :one
SELECT COUNT(*) FROM group_members WHERE user_id = $1;
