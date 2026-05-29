-- name: GetExpenseCategoryByID :one
SELECT * FROM expense_categories
WHERE id = $1;

-- name: ListExpenseCategoriesByGroup :many
SELECT * FROM expense_categories
WHERE group_id = $1
ORDER BY name;

-- name: CreateExpenseCategory :one
INSERT INTO expense_categories (group_id, name)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateExpenseCategory :one
UPDATE expense_categories
SET name = $3
WHERE id = $1 AND group_id = $2
RETURNING *;

-- name: DeleteExpenseCategory :exec
DELETE FROM expense_categories
WHERE id = $1 AND group_id = $2;
