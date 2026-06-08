-- name: GetExpenseByID :one
SELECT * FROM expenses
WHERE id = $1;

-- name: ListExpensesByGroup :many
SELECT * FROM expenses
WHERE group_id = $1
ORDER BY competence_date DESC, created_at DESC;

-- name: ListExpensesByUser :many
SELECT e.* FROM expenses e
JOIN expense_splits es ON es.expense_id = e.id
WHERE es.user_id = $1
ORDER BY e.competence_date DESC, e.created_at DESC;

-- name: ListExpensesByCategory :many
SELECT * FROM expenses
WHERE category_id = $1
ORDER BY competence_date DESC;

-- name: CreateExpense :one
INSERT INTO expenses (group_id, description, amount, category_id, competence_date, due_date, paid_by, split_mode, installments, first_due_date, is_fixed, parent_expense_id, installment_index, installment_total)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: UpdateExpense :one
UPDATE expenses
SET description = COALESCE($2, description),
    amount = COALESCE($3, amount),
    competence_date = COALESCE($4, competence_date),
    due_date = COALESCE($5, due_date),
    category_id = COALESCE($6, category_id),
    split_mode = COALESCE($7, split_mode),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteExpense :exec
DELETE FROM expenses
WHERE id = $1;

-- name: GetTotalExpensesByGroup :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC(12,2) FROM expenses
WHERE group_id = $1;

-- name: ListExpensesByParent :many
SELECT * FROM expenses
WHERE parent_expense_id = $1
ORDER BY installment_index;

-- name: DeleteExpensesByParent :exec
DELETE FROM expenses
WHERE parent_expense_id = $1;
