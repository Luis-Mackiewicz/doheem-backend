-- name: GetExpenseByID :one
SELECT * FROM expenses
WHERE id = $1;

-- name: ListExpensesByGroup :many
SELECT * FROM expenses
WHERE group_id = $1
ORDER BY expense_date DESC, created_at DESC;

-- name: ListExpensesByUser :many
SELECT e.* FROM expenses e
JOIN expense_splits es ON es.expense_id = e.id
WHERE es.user_id = $1
ORDER BY e.expense_date DESC, e.created_at DESC;

-- name: ListExpensesByCategory :many
SELECT * FROM expenses
WHERE category_id = $1
ORDER BY expense_date DESC;

-- name: CreateExpense :one
INSERT INTO expenses (group_id, created_by, description, total_amount, expense_date, due_date, category_id, split_type, is_installment, installment_count)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: UpdateExpense :one
UPDATE expenses
SET description = COALESCE($2, description),
    total_amount = COALESCE($3, total_amount),
    expense_date = COALESCE($4, expense_date),
    due_date = $5,
    category_id = $6,
    split_type = COALESCE($7, split_type),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteExpense :exec
DELETE FROM expenses
WHERE id = $1;

-- name: GetTotalExpensesByGroup :one
SELECT COALESCE(SUM(total_amount), 0)::NUMERIC(12,2) FROM expenses
WHERE group_id = $1;
