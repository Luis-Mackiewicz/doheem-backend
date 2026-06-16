package notification

import (
	"context"

	"doheem-backend/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type NotificationRepo struct {
	q *db.Queries
}

func NewNotificationRepo(q *db.Queries) *NotificationRepo {
	return &NotificationRepo{q: q}
}

func (r *NotificationRepo) GetByID(ctx context.Context, id string) (Notification, error) {
	n, err := r.q.GetNotificationByID(ctx, db.UUIDFromString(id))
	if err != nil {
		return Notification{}, err
	}
	return toNotification(n), nil
}

func (r *NotificationRepo) ListByUser(ctx context.Context, userID string, limit, offset int32) ([]Notification, error) {
	notifications, err := r.q.ListNotificationsByUser(ctx, db.ListNotificationsByUserParams{
		UserID: db.UUIDFromString(userID),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return toNotifications(notifications), nil
}

func (r *NotificationRepo) ListByUserSearch(ctx context.Context, userID, search string, limit, offset int32) ([]Notification, error) {
	notifications, err := r.q.ListNotificationsByUserSearch(ctx, db.ListNotificationsByUserSearchParams{
		UserID:  db.UUIDFromString(userID),
		Column2: pgtype.Text{String: search, Valid: true},
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}
	return toNotifications(notifications), nil
}

func (r *NotificationRepo) ListUnreadByUser(ctx context.Context, userID string) ([]Notification, error) {
	notifications, err := r.q.ListUnreadNotificationsByUser(ctx, db.UUIDFromString(userID))
	if err != nil {
		return nil, err
	}
	return toNotifications(notifications), nil
}

func (r *NotificationRepo) Create(ctx context.Context, params CreateNotificationParams) (Notification, error) {
	n, err := r.q.CreateNotification(ctx, db.CreateNotificationParams{
		UserID:    db.UUIDFromString(params.UserID),
		Type:      params.Type,
		Title:     params.Title,
		Message:   params.Message,
		RelatedID: db.UUIDFromStringPtr(params.RelatedID),
	})
	if err != nil {
		return Notification{}, err
	}
	return toNotification(n), nil
}

func (r *NotificationRepo) MarkAsRead(ctx context.Context, id, userID string) error {
	return r.q.MarkNotificationAsRead(ctx, db.MarkNotificationAsReadParams{
		ID:     db.UUIDFromString(id),
		UserID: db.UUIDFromString(userID),
	})
}

func (r *NotificationRepo) MarkAllAsRead(ctx context.Context, userID string) error {
	return r.q.MarkAllNotificationsAsRead(ctx, db.UUIDFromString(userID))
}

func (r *NotificationRepo) CountUnread(ctx context.Context, userID string) (int64, error) {
	return r.q.CountUnreadNotifications(ctx, db.UUIDFromString(userID))
}

func (r *NotificationRepo) CountByUser(ctx context.Context, userID string) (int64, error) {
	return r.q.CountNotificationsByUser(ctx, db.UUIDFromString(userID))
}

func (r *NotificationRepo) Delete(ctx context.Context, id, userID string) error {
	return r.q.DeleteNotification(ctx, db.DeleteNotificationParams{
		ID:     db.UUIDFromString(id),
		UserID: db.UUIDFromString(userID),
	})
}

func (r *NotificationRepo) DeleteAll(ctx context.Context, userID string) error {
	return r.q.DeleteAllNotificationsByUser(ctx, db.UUIDFromString(userID))
}

func toNotification(n db.Notification) Notification {
	return Notification{
		ID:        db.UUIDToString(n.ID),
		UserID:    db.UUIDToString(n.UserID),
		Type:      n.Type,
		Title:     n.Title,
		Message:   n.Message,
		IsRead:    n.IsRead,
		RelatedID: db.UUIDToStringPtr(n.RelatedID),
		CreatedAt: n.CreatedAt.Time,
	}
}

func toNotifications(notifications []db.Notification) []Notification {
	result := make([]Notification, len(notifications))
	for i, n := range notifications {
		result[i] = toNotification(n)
	}
	return result
}
