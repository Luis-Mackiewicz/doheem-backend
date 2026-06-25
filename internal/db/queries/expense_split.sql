-- name: GetExpenseSplitByID :one
SELECT * FROM expense_splits
WHERE id = $1;

-- name: ListExpenseSplitsByExpense :many
SELECT es.*, u.name AS user_name, u.email AS user_email
FROM expense_splits es
JOIN users u ON u.id = es.user_id
WHERE es.expense_id = $1
ORDER BY es.amount DESC;

-- name: ListExpenseSplitsByExpenseIDs :many
SELECT es.*, u.name AS user_name, u.email AS user_email
FROM expense_splits es
JOIN users u ON u.id = es.user_id
WHERE es.expense_id = ANY($1::uuid[])
ORDER BY es.expense_id, es.amount DESC;

-- name: ListExpenseSplitsByUser :many
SELECT es.*, e.description AS expense_description
FROM expense_splits es
JOIN expenses e ON e.id = es.expense_id
WHERE es.user_id = $1 AND e.group_id = $2
ORDER BY e.competence_date DESC;

-- name: CreateExpenseSplit :one
INSERT INTO expense_splits (expense_id, user_id, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateExpenseSplits :copyfrom
INSERT INTO expense_splits (expense_id, user_id, amount)
VALUES ($1, $2, $3);

-- name: MarkExpenseSplitAsPaid :exec
UPDATE expense_splits
SET is_paid = true,
    paid_at = NOW(),
    receipt_data = $2,
    receipt_type = $3,
    receipt_file_name = $4
WHERE id = $1;

-- name: DeleteExpenseSplitsByExpense :exec
DELETE FROM expense_splits
WHERE expense_id = $1;

-- name: HasExpensePaidSplits :one
SELECT EXISTS(
  SELECT 1 FROM expense_splits
  WHERE expense_id = $1 AND is_paid = true
);

-- name: GetUserBalanceInGroup :one
SELECT COALESCE(SUM(es.amount), 0)::NUMERIC(12,2) AS total_owed,
       COALESCE(SUM(CASE WHEN es.is_paid THEN es.amount ELSE 0 END), 0)::NUMERIC(12,2) AS total_paid
FROM expense_splits es
JOIN expenses e ON e.id = es.expense_id
WHERE es.user_id = $1 AND e.group_id = $2;
