-- name: GetInstallmentByID :one
SELECT * FROM installments
WHERE id = $1;

-- name: ListInstallmentsByExpense :many
SELECT * FROM installments
WHERE expense_id = $1
ORDER BY installment_number;

-- name: CreateInstallment :one
INSERT INTO installments (expense_id, installment_number, amount, due_date)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CreateInstallments :copyfrom
INSERT INTO installments (expense_id, installment_number, amount, due_date)
VALUES ($1, $2, $3, $4);

-- name: MarkInstallmentAsPaid :exec
UPDATE installments
SET is_paid = true,
    paid_at = NOW()
WHERE id = $1;

-- name: DeleteInstallmentsByExpense :exec
DELETE FROM installments
WHERE expense_id = $1;
