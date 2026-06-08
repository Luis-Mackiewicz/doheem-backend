-- name: GetAuditLogByID :one
SELECT * FROM audit_logs
WHERE id = $1;

-- name: ListAuditLogsByGroup :many
SELECT al.*, u.name AS user_name
FROM audit_logs al
LEFT JOIN users u ON u.id = al.user_id
WHERE al.group_id = $1
ORDER BY al.created_at DESC;

-- name: ListAuditLogsByEntity :many
SELECT * FROM audit_logs
WHERE entity_type = $1 AND entity_id = $2
ORDER BY created_at DESC;

-- name: CreateAuditLog :one
INSERT INTO audit_logs (group_id, user_id, entity_type, entity_id, action, changes)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
