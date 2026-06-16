package notification

import (
	"context"
	"time"
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

func (s *NotificationService) ListByUser(ctx context.Context, userID, search string, limit, offset int32) ([]Notification, int64, error) {
	if search != "" {
		notifications, err := s.repo.ListByUserSearch(ctx, userID, search, limit, offset)
		if err != nil {
			return nil, 0, err
		}
		return notifications, 0, nil
	}
	notifications, err := s.repo.ListByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.CountByUser(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	return notifications, total, nil
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

func (s *NotificationService) DeleteAll(ctx context.Context, userID string) error {
	return s.repo.DeleteAll(ctx, userID)
}

func (s *NotificationService) CanSendReminder(ctx context.Context, userID, expenseID string) error {
	lastReminders, err := s.repo.ListByUserSearch(ctx, userID, expenseID, MaxReminders, 0)
	if err != nil {
		return err
	}
	if len(lastReminders) >= MaxReminders {
		return ErrReminderLimitExceeded
	}
	if len(lastReminders) > 0 {
		last := lastReminders[0].CreatedAt
		diff := time.Since(last).Hours() / 24
		if diff < MinIntervalDays {
			return ErrReminderTooSoon
		}
	}
	return nil
}
