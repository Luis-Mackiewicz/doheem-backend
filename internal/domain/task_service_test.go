package domain

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTaskRepo struct {
	mock.Mock
}

func (m *mockTaskRepo) GetByID(ctx context.Context, id string) (Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Task), args.Error(1)
}

func (m *mockTaskRepo) ListByGroup(ctx context.Context, groupID string) ([]Task, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).([]Task), args.Error(1)
}

func (m *mockTaskRepo) ListByAssignee(ctx context.Context, userID, groupID string) ([]Task, error) {
	args := m.Called(ctx, userID, groupID)
	return args.Get(0).([]Task), args.Error(1)
}

func (m *mockTaskRepo) Create(ctx context.Context, params CreateTaskParams) (Task, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(Task), args.Error(1)
}

func (m *mockTaskRepo) Update(ctx context.Context, id string, params UpdateTaskParams) (Task, error) {
	args := m.Called(ctx, id, params)
	return args.Get(0).(Task), args.Error(1)
}

func (m *mockTaskRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockTaskOccurrenceRepo struct {
	mock.Mock
}

func (m *mockTaskOccurrenceRepo) GetByID(ctx context.Context, id string) (TaskOccurrence, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(TaskOccurrence), args.Error(1)
}

func (m *mockTaskOccurrenceRepo) ListByTask(ctx context.Context, taskID string) ([]TaskOccurrence, error) {
	args := m.Called(ctx, taskID)
	return args.Get(0).([]TaskOccurrence), args.Error(1)
}

func (m *mockTaskOccurrenceRepo) ListPendingByUser(ctx context.Context, userID string) ([]TaskOccurrenceWithTask, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]TaskOccurrenceWithTask), args.Error(1)
}

func (m *mockTaskOccurrenceRepo) ListByDateRange(ctx context.Context, groupID string, start, end time.Time) ([]TaskOccurrenceWithTask, error) {
	args := m.Called(ctx, groupID, start, end)
	return args.Get(0).([]TaskOccurrenceWithTask), args.Error(1)
}

func (m *mockTaskOccurrenceRepo) Create(ctx context.Context, taskID string, dueDate time.Time, status string) (TaskOccurrence, error) {
	args := m.Called(ctx, taskID, dueDate, status)
	return args.Get(0).(TaskOccurrence), args.Error(1)
}

func (m *mockTaskOccurrenceRepo) Complete(ctx context.Context, id string, completedBy string) error {
	args := m.Called(ctx, id, completedBy)
	return args.Error(0)
}

func (m *mockTaskOccurrenceRepo) Discard(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockTaskOccurrenceRepo) MarkAsOverdue(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockTaskOccurrenceRepo) DeleteByTask(ctx context.Context, taskID string) error {
	args := m.Called(ctx, taskID)
	return args.Error(0)
}

func TestTaskService_GetByID_NotFound(t *testing.T) {
	mockTask := new(mockTaskRepo)
	svc := NewTaskService(mockTask, new(mockTaskOccurrenceRepo))
	ctx := context.Background()

	mockTask.On("GetByID", ctx, "999").Return(Task{}, assert.AnError)

	_, err := svc.GetByID(ctx, "999")

	assert.ErrorIs(t, err, ErrTaskNotFound)
}

func TestTaskService_GetByID_Success(t *testing.T) {
	mockTask := new(mockTaskRepo)
	svc := NewTaskService(mockTask, new(mockTaskOccurrenceRepo))
	ctx := context.Background()

	mockTask.On("GetByID", ctx, "1").Return(Task{ID: "1", Title: "Test Task"}, nil)

	task, err := svc.GetByID(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, "1", task.ID)
}

func TestTaskService_CompleteOccurrence(t *testing.T) {
	mockOccurrence := new(mockTaskOccurrenceRepo)
	svc := NewTaskService(new(mockTaskRepo), mockOccurrence)
	ctx := context.Background()

	mockOccurrence.On("Complete", ctx, "o1", "u1").Return(nil)

	err := svc.CompleteOccurrence(ctx, "o1", "u1")

	assert.NoError(t, err)
}

func TestTaskService_DiscardOccurrence(t *testing.T) {
	mockOccurrence := new(mockTaskOccurrenceRepo)
	svc := NewTaskService(new(mockTaskRepo), mockOccurrence)
	ctx := context.Background()

	mockOccurrence.On("Discard", ctx, "o1").Return(nil)

	err := svc.DiscardOccurrence(ctx, "o1")

	assert.NoError(t, err)
}
