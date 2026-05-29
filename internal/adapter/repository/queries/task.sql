-- name: GetTaskByID :one
SELECT * FROM tasks
WHERE id = $1;

-- name: ListTasksByGroup :many
SELECT * FROM tasks
WHERE group_id = $1
ORDER BY created_at DESC;

-- name: ListTasksByAssignee :many
SELECT * FROM tasks
WHERE assigned_to = $1 AND group_id = $2
ORDER BY created_at DESC;

-- name: CreateTask :one
INSERT INTO tasks (group_id, title, description, assigned_to, category, is_recurring, recurring_period, recurring_ended_at, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateTask :one
UPDATE tasks
SET title = COALESCE($2, title),
    description = $3,
    assigned_to = $4,
    category = $5,
    is_recurring = COALESCE($6, is_recurring),
    recurring_period = $7,
    recurring_ended_at = $8,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;
