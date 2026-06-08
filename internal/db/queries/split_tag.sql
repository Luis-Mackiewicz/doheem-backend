-- name: GetSplitTagByID :one
SELECT * FROM split_tags
WHERE id = $1;

-- name: ListSplitTagsByGroup :many
SELECT * FROM split_tags
WHERE group_id = $1
ORDER BY name;

-- name: CreateSplitTag :one
INSERT INTO split_tags (group_id, name, created_by)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteSplitTag :exec
DELETE FROM split_tags
WHERE id = $1 AND group_id = $2;

-- name: ListSplitTagMembers :many
SELECT stm.*, u.name AS user_name
FROM split_tag_members stm
JOIN users u ON u.id = stm.user_id
WHERE stm.split_tag_id = $1;

-- name: AddSplitTagMember :one
INSERT INTO split_tag_members (split_tag_id, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: RemoveSplitTagMember :exec
DELETE FROM split_tag_members
WHERE split_tag_id = $1 AND user_id = $2;
