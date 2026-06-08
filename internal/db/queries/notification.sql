-- name: GetNotificationByID :one
SELECT * FROM notifications
WHERE id = $1;

-- name: ListNotificationsByUser :many
SELECT * FROM notifications
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListUnreadNotificationsByUser :many
SELECT * FROM notifications
WHERE user_id = $1 AND is_read = false
ORDER BY created_at DESC;

-- name: CreateNotification :one
INSERT INTO notifications (user_id, type, title, message, related_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: MarkNotificationAsRead :exec
UPDATE notifications
SET is_read = true
WHERE id = $1 AND user_id = $2;

-- name: MarkAllNotificationsAsRead :exec
UPDATE notifications
SET is_read = true
WHERE user_id = $1 AND is_read = false;

-- name: CountUnreadNotifications :one
SELECT COUNT(*) FROM notifications
WHERE user_id = $1 AND is_read = false;

-- name: DeleteNotification :exec
DELETE FROM notifications
WHERE id = $1 AND user_id = $2;
