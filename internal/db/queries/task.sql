-- name: GetTaskByID :one
SELECT * FROM tasks
WHERE id = $1;

-- name: ListTasksByGroup :many
SELECT * FROM tasks
WHERE group_id = $1
ORDER BY position, created_at DESC;

-- name: ListTasksByAssignee :many
SELECT * FROM tasks
WHERE assigned_to = $1 AND group_id = $2
ORDER BY position, created_at DESC;

-- name: CreateTask :one
INSERT INTO tasks (group_id, title, description, assigned_to, created_by, due_date)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateTask :one
UPDATE tasks
SET title = COALESCE($2, title),
    description = COALESCE($3, description),
    assigned_to = COALESCE($4, assigned_to),
    status = COALESCE($5, status),
    position = COALESCE($6, position),
    due_date = COALESCE($7, due_date),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;
