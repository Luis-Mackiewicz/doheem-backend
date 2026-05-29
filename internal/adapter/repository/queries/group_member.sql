-- name: GetGroupMemberByID :one
SELECT * FROM group_members
WHERE id = $1;

-- name: GetGroupMember :one
SELECT * FROM group_members
WHERE group_id = $1 AND user_id = $2;

-- name: ListGroupMembers :many
SELECT gm.*, u.name, u.email, u.avatar_url
FROM group_members gm
JOIN users u ON u.id = gm.user_id
WHERE gm.group_id = $1 AND gm.is_active = true
ORDER BY gm.joined_at;

-- name: CreateGroupMember :one
INSERT INTO group_members (group_id, user_id, role)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateGroupMemberRole :one
UPDATE group_members
SET role = $3
WHERE group_id = $1 AND user_id = $2
RETURNING *;

-- name: RemoveGroupMember :exec
UPDATE group_members
SET is_active = false,
    left_at = NOW()
WHERE group_id = $1 AND user_id = $2;

-- name: CountActiveGroupMembers :one
SELECT COUNT(*) FROM group_members
WHERE group_id = $1 AND is_active = true;
