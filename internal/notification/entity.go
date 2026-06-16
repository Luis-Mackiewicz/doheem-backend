package notification

import (
	"context"
	"errors"
	"time"
)

type Notification struct {
	ID        string
	UserID    string
	Type      string
	Title     string
	Message   string
	IsRead    bool
	RelatedID *string
	CreatedAt time.Time
}

type CreateNotificationParams struct {
	UserID    string
	Type      string
	Title     string
	Message   string
	RelatedID *string
}

type NotificationRepository interface {
	GetByID(ctx context.Context, id string) (Notification, error)
	ListByUser(ctx context.Context, userID string, limit, offset int32) ([]Notification, error)
	ListByUserSearch(ctx context.Context, userID, search string, limit, offset int32) ([]Notification, error)
	ListUnreadByUser(ctx context.Context, userID string) ([]Notification, error)
	Create(ctx context.Context, params CreateNotificationParams) (Notification, error)
	MarkAsRead(ctx context.Context, id, userID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	CountUnread(ctx context.Context, userID string) (int64, error)
	CountByUser(ctx context.Context, userID string) (int64, error)
	Delete(ctx context.Context, id, userID string) error
	DeleteAll(ctx context.Context, userID string) error
}

const (
	MaxReminders      = 5
	MinIntervalDays   = 3
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrReminderLimitExceeded = errors.New("limite de lembretes atingido para esta despesa")
	ErrReminderTooSoon       = errors.New("aguarde o intervalo mínimo entre lembretes")
)
