package notification

import (
	"context"
)

type NotificationService struct {
	repo NotificationRepository
}

func NewNotificationService(repo NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) Create(ctx context.Context, params CreateNotificationParams) (Notification, error) {
	return s.repo.Create(ctx, params)
}

func (s *NotificationService) GetByID(ctx context.Context, id string) (Notification, error) {
	notification, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Notification{}, ErrNotificationNotFound
	}
	return notification, nil
}

func (s *NotificationService) ListByUser(ctx context.Context, userID string, limit, offset int32) ([]Notification, error) {
	return s.repo.ListByUser(ctx, userID, limit, offset)
}

func (s *NotificationService) ListUnreadByUser(ctx context.Context, userID string) ([]Notification, error) {
	return s.repo.ListUnreadByUser(ctx, userID)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, id, userID string) error {
	return s.repo.MarkAsRead(ctx, id, userID)
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

func (s *NotificationService) CountUnread(ctx context.Context, userID string) (int64, error) {
	return s.repo.CountUnread(ctx, userID)
}

func (s *NotificationService) Delete(ctx context.Context, id, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}
