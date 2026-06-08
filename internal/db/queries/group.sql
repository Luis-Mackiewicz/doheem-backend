-- name: GetGroupByID :one
SELECT * FROM groups
WHERE id = $1;

-- name: ListGroupsByUserID :many
SELECT g.* FROM groups g
JOIN group_members gm ON gm.group_id = g.id
WHERE gm.user_id = $1
ORDER BY g.name;

-- name: CreateGroup :one
INSERT INTO groups (name, invite_token)
VALUES ($1, gen_random_uuid()::text)
RETURNING *;

-- name: UpdateGroup :one
UPDATE groups
SET name = COALESCE($2, name),
    description = COALESCE($3, description),
    monthly_fee = COALESCE($4, monthly_fee),
    cnpj = COALESCE($5, cnpj),
    cep = COALESCE($6, cep),
    photo_url = COALESCE($7, photo_url),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: RegenerateInviteToken :one
UPDATE groups
SET invite_token = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
