package repository

import (
	"context"

	"doheem-backend/internal/adapter/repository/db"
	"doheem-backend/internal/domain"
)

type NotificationRepo struct {
	q *db.Queries
}

func NewNotificationRepo(q *db.Queries) *NotificationRepo {
	return &NotificationRepo{q: q}
}

func (r *NotificationRepo) GetByID(ctx context.Context, id string) (domain.Notification, error) {
	n, err := r.q.GetNotificationByID(ctx, uuidFromString(id))
	if err != nil {
		return domain.Notification{}, err
	}
	return domainNotification(n), nil
}

func (r *NotificationRepo) ListByUser(ctx context.Context, userID string, limit, offset int32) ([]domain.Notification, error) {
	notifications, err := r.q.ListNotificationsByUser(ctx, db.ListNotificationsByUserParams{
		UserID: uuidFromString(userID),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return domainNotifications(notifications), nil
}

func (r *NotificationRepo) ListUnreadByUser(ctx context.Context, userID string) ([]domain.Notification, error) {
	notifications, err := r.q.ListUnreadNotificationsByUser(ctx, uuidFromString(userID))
	if err != nil {
		return nil, err
	}
	return domainNotifications(notifications), nil
}

func (r *NotificationRepo) Create(ctx context.Context, params domain.CreateNotificationParams) (domain.Notification, error) {
	n, err := r.q.CreateNotification(ctx, db.CreateNotificationParams{
		UserID:  uuidFromString(params.UserID),
		GroupID: uuidFromStringPtr(params.GroupID),
		Type:    params.Type,
		Title:   params.Title,
		Message: params.Message,
	})
	if err != nil {
		return domain.Notification{}, err
	}
	return domainNotification(n), nil
}

func (r *NotificationRepo) MarkAsRead(ctx context.Context, id, userID string) error {
	return r.q.MarkNotificationAsRead(ctx, db.MarkNotificationAsReadParams{
		ID:     uuidFromString(id),
		UserID: uuidFromString(userID),
	})
}

func (r *NotificationRepo) MarkAllAsRead(ctx context.Context, userID string) error {
	return r.q.MarkAllNotificationsAsRead(ctx, uuidFromString(userID))
}

func (r *NotificationRepo) CountUnread(ctx context.Context, userID string) (int64, error) {
	return r.q.CountUnreadNotifications(ctx, uuidFromString(userID))
}

func (r *NotificationRepo) Delete(ctx context.Context, id, userID string) error {
	return r.q.DeleteNotification(ctx, db.DeleteNotificationParams{
		ID:     uuidFromString(id),
		UserID: uuidFromString(userID),
	})
}
