package notification

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockNotificationRepo struct {
	mock.Mock
}

func (m *mockNotificationRepo) GetByID(ctx context.Context, id string) (Notification, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Notification), args.Error(1)
}

func (m *mockNotificationRepo) ListByUser(ctx context.Context, userID string, limit, offset int32) ([]Notification, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]Notification), args.Error(1)
}

func (m *mockNotificationRepo) ListUnreadByUser(ctx context.Context, userID string) ([]Notification, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]Notification), args.Error(1)
}

func (m *mockNotificationRepo) Create(ctx context.Context, params CreateNotificationParams) (Notification, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(Notification), args.Error(1)
}

func (m *mockNotificationRepo) MarkAsRead(ctx context.Context, id, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *mockNotificationRepo) MarkAllAsRead(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockNotificationRepo) CountUnread(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockNotificationRepo) Delete(ctx context.Context, id, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *mockNotificationRepo) ListByUserSearch(ctx context.Context, userID, search string, limit, offset int32) ([]Notification, error) {
	args := m.Called(ctx, userID, search, limit, offset)
	return args.Get(0).([]Notification), args.Error(1)
}

func (m *mockNotificationRepo) CountByUser(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockNotificationRepo) DeleteAll(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestNotificationService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(mockNotificationRepo)
	svc := NewNotificationService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, "999").Return(Notification{}, assert.AnError)

	_, err := svc.GetByID(ctx, "999")

	assert.ErrorIs(t, err, ErrNotificationNotFound)
}

func TestNotificationService_GetByID_Success(t *testing.T) {
	mockRepo := new(mockNotificationRepo)
	svc := NewNotificationService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, "1").Return(Notification{ID: "1", Title: "Test"}, nil)

	n, err := svc.GetByID(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, "1", n.ID)
}

func TestNotificationService_Create(t *testing.T) {
	mockRepo := new(mockNotificationRepo)
	svc := NewNotificationService(mockRepo)
	ctx := context.Background()

	params := CreateNotificationParams{UserID: "u1", Title: "Hello", Message: "World"}
	mockRepo.On("Create", ctx, params).Return(Notification{ID: "n1", Title: "Hello"}, nil)

	n, err := svc.Create(ctx, params)

	assert.NoError(t, err)
	assert.Equal(t, "n1", n.ID)
}

func TestNotificationService_MarkAllAsRead(t *testing.T) {
	mockRepo := new(mockNotificationRepo)
	svc := NewNotificationService(mockRepo)
	ctx := context.Background()

	mockRepo.On("MarkAllAsRead", ctx, "u1").Return(nil)

	err := svc.MarkAllAsRead(ctx, "u1")

	assert.NoError(t, err)
}

func TestNotificationService_CountUnread(t *testing.T) {
	mockRepo := new(mockNotificationRepo)
	svc := NewNotificationService(mockRepo)
	ctx := context.Background()

	mockRepo.On("CountUnread", ctx, "u1").Return(int64(5), nil)

	count, err := svc.CountUnread(ctx, "u1")

	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
}
