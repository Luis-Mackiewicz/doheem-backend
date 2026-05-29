-- name: GetInviteByID :one
SELECT * FROM invites
WHERE id = $1;

-- name: GetInviteByCode :one
SELECT i.*, g.name AS group_name
FROM invites i
JOIN groups g ON g.id = i.group_id
WHERE i.code = $1 AND i.used_at IS NULL AND i.revoked_at IS NULL AND i.expires_at > NOW();

-- name: ListInvitesByGroup :many
SELECT i.*, u.name AS created_by_name
FROM invites i
JOIN users u ON u.id = i.created_by
WHERE i.group_id = $1
ORDER BY i.created_at DESC;

-- name: ListPendingInvitesByUser :many
SELECT i.*, g.name AS group_name
FROM invites i
JOIN groups g ON g.id = i.group_id
JOIN group_members gm ON gm.group_id = i.group_id AND gm.user_id = $1
WHERE i.used_at IS NULL AND i.revoked_at IS NULL AND i.expires_at > NOW()
ORDER BY i.created_at DESC;

-- name: CreateInvite :one
INSERT INTO invites (group_id, code, created_by, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UseInvite :exec
UPDATE invites
SET used_at = NOW()
WHERE id = $1 AND used_at IS NULL;

-- name: RevokeInvite :exec
UPDATE invites
SET revoked_at = NOW()
WHERE id = $1 AND used_at IS NULL;
