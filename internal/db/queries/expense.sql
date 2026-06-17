-- name: GetExpenseByID :one
SELECT * FROM expenses
WHERE id = $1;

-- name: ListExpensesByGroup :many
SELECT * FROM expenses
WHERE group_id = $1
  AND ($2::DATE IS NULL OR competence_date >= $2)
  AND ($3::DATE IS NULL OR competence_date <= $3)
  AND (installments = 1 OR parent_expense_id IS NOT NULL)
ORDER BY competence_date DESC, created_at DESC
LIMIT $4 OFFSET $5;

-- name: ListExpensesByGroupFull :many
SELECT e.*, COUNT(*) OVER() AS total_count
FROM expenses e
WHERE e.group_id = $1
  AND ($2::DATE IS NULL OR e.competence_date >= $2)
  AND ($3::DATE IS NULL OR e.competence_date <= $3)
  AND ($4::TEXT IS NULL OR $4 = '' OR e.description ILIKE '%' || $4 || '%')
  AND ($5::UUID IS NULL OR EXISTS (
    SELECT 1 FROM expense_splits es WHERE es.expense_id = e.id AND es.user_id = $5
  ))
  AND (e.installments = 1 OR e.parent_expense_id IS NOT NULL)
ORDER BY e.competence_date DESC, e.created_at DESC
LIMIT $6 OFFSET $7;

-- name: CountExpensesByGroup :one
SELECT COUNT(*) FROM expenses
WHERE group_id = $1
  AND ($2::DATE IS NULL OR competence_date >= $2)
  AND ($3::DATE IS NULL OR competence_date <= $3)
  AND (installments = 1 OR parent_expense_id IS NOT NULL);

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
INSERT INTO expenses (group_id, description, amount, category_id, competence_date, due_date, paid_by, split_mode, installments, first_due_date, is_fixed, parent_expense_id, installment_index, installment_total, created_by, fixed_source_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
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
WHERE group_id = $1
  AND (installments = 1 OR parent_expense_id IS NOT NULL);

-- name: ListExpensesByParent :many
SELECT * FROM expenses
WHERE parent_expense_id = $1
ORDER BY installment_index;

-- name: DeleteExpensesByParent :exec
DELETE FROM expenses
WHERE parent_expense_id = $1;

-- name: ListFixedOrigins :many
SELECT * FROM expenses
WHERE is_fixed = true
  AND installments = 1
  AND fixed_source_id IS NULL;

-- name: CountCloneByMonth :one
SELECT COUNT(*) FROM expenses
WHERE fixed_source_id = $1
  AND competence_date >= $2
  AND competence_date <= $3;
