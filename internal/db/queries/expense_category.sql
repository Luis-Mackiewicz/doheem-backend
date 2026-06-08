-- name: GetExpenseCategoryByID :one
SELECT * FROM expense_categories
WHERE id = $1;

-- name: ListExpenseCategories :many
SELECT * FROM expense_categories
ORDER BY label;

-- name: CreateExpenseCategory :one
INSERT INTO expense_categories (slug, label)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateExpenseCategory :one
UPDATE expense_categories
SET slug = $2,
    label = $3
WHERE id = $1
RETURNING *;

-- name: DeleteExpenseCategory :exec
DELETE FROM expense_categories
WHERE id = $1;
