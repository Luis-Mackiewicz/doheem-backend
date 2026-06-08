-- name: GetPaymentAttachmentByID :one
SELECT * FROM payment_attachments
WHERE id = $1;

-- name: ListPaymentAttachmentsByPayment :many
SELECT * FROM payment_attachments
WHERE payment_id = $1
ORDER BY uploaded_at;

-- name: CreatePaymentAttachment :one
INSERT INTO payment_attachments (payment_id, file_path, file_type, file_size)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeletePaymentAttachment :exec
DELETE FROM payment_attachments
WHERE id = $1;
