package domain

import (
	"context"
	"time"
)

type Notification struct {
	ID        string
	UserID    string
	GroupID   *string
	Type      string
	Title     string
	Message   string
	IsRead    bool
	ReadAt    *time.Time
	CreatedAt time.Time
}

type CreateNotificationParams struct {
	UserID  string
	GroupID *string
	Type    string
	Title   string
	Message string
}

type NotificationRepository interface {
	GetByID(ctx context.Context, id string) (Notification, error)
	ListByUser(ctx context.Context, userID string, limit, offset int32) ([]Notification, error)
	ListUnreadByUser(ctx context.Context, userID string) ([]Notification, error)
	Create(ctx context.Context, params CreateNotificationParams) (Notification, error)
	MarkAsRead(ctx context.Context, id, userID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	CountUnread(ctx context.Context, userID string) (int64, error)
	Delete(ctx context.Context, id, userID string) error
}
