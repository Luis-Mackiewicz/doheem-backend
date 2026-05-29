-- name: GetTaskOccurrenceByID :one
SELECT * FROM task_occurrences
WHERE id = $1;

-- name: ListTaskOccurrencesByTask :many
SELECT * FROM task_occurrences
WHERE task_id = $1
ORDER BY due_date;

-- name: ListPendingTaskOccurrencesByUser :many
SELECT to2.*, t.title AS task_title, t.group_id
FROM task_occurrences to2
JOIN tasks t ON t.id = to2.task_id
WHERE to2.status = 'pending' AND t.assigned_to = $1
ORDER BY to2.due_date;

-- name: ListTaskOccurrencesByDateRange :many
SELECT to2.*, t.title AS task_title, t.group_id
FROM task_occurrences to2
JOIN tasks t ON t.id = to2.task_id
WHERE t.group_id = $1 AND to2.due_date BETWEEN $2 AND $3
ORDER BY to2.due_date;

-- name: CreateTaskOccurrence :one
INSERT INTO task_occurrences (task_id, due_date, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CompleteTaskOccurrence :exec
UPDATE task_occurrences
SET status = 'completed',
    completed_by = $2,
    completed_at = NOW()
WHERE id = $1 AND status = 'pending';

-- name: DiscardTaskOccurrence :exec
UPDATE task_occurrences
SET status = 'discarded',
    discarded_at = NOW()
WHERE id = $1 AND status = 'pending';

-- name: MarkTaskOccurrenceAsOverdue :exec
UPDATE task_occurrences
SET status = 'overdue'
WHERE id = $1 AND status = 'pending' AND due_date < CURRENT_DATE;

-- name: DeleteTaskOccurrencesByTask :exec
DELETE FROM task_occurrences
WHERE task_id = $1;
