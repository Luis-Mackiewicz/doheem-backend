-- name: GetPaymentByID :one
SELECT * FROM payments
WHERE id = $1;

-- name: ListPaymentsByGroup :many
SELECT p.*, pu.name AS payer_name, ru.name AS receiver_name
FROM payments p
JOIN users pu ON pu.id = p.payer_id
JOIN users ru ON ru.id = p.receiver_id
WHERE p.group_id = $1
ORDER BY p.payment_date DESC, p.created_at DESC;

-- name: ListPaymentsByUser :many
SELECT p.*, g.name AS group_name
FROM payments p
JOIN groups g ON g.id = p.group_id
WHERE p.payer_id = $1 OR p.receiver_id = $1
ORDER BY p.payment_date DESC, p.created_at DESC;

-- name: ListPendingPaymentsByUser :many
SELECT * FROM payments
WHERE (payer_id = $1 OR receiver_id = $1)
  AND status = 'pending'
ORDER BY payment_date;

-- name: CreatePayment :one
INSERT INTO payments (group_id, payer_id, receiver_id, amount, payment_date, notes)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ConfirmPayment :exec
UPDATE payments
SET status = 'confirmed',
    confirmed_at = NOW()
WHERE id = $1 AND status = 'pending';

-- name: CancelPayment :exec
UPDATE payments
SET status = 'cancelled',
    cancelled_at = NOW()
WHERE id = $1 AND status = 'pending';

-- name: DeletePayment :exec
DELETE FROM payments
WHERE id = $1;
