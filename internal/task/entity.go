package task

import (
	"context"
	"errors"
	"time"
)

type Task struct {
	ID               string
	GroupID          string
	Title            string
	Description      *string
	AssignedTo       *string
	Category         *string
	IsRecurring      bool
	RecurringPeriod  *string
	RecurringEndedAt *time.Time
	CreatedBy        string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CreateTaskParams struct {
	GroupID          string
	Title            string
	Description      *string
	AssignedTo       *string
	Category         *string
	IsRecurring      bool
	RecurringPeriod  *string
	RecurringEndedAt *time.Time
	CreatedBy        string
}

type UpdateTaskParams struct {
	Title            *string
	Description      *string
	AssignedTo       *string
	Category         *string
	IsRecurring      *bool
	RecurringPeriod  *string
	RecurringEndedAt *time.Time
}

type TaskOccurrence struct {
	ID          string
	TaskID      string
	DueDate     time.Time
	Status      string
	CompletedBy *string
	CompletedAt *time.Time
	DiscardedAt *time.Time
	CreatedAt   time.Time
}

type TaskOccurrenceWithTask struct {
	TaskOccurrence
	TaskTitle string
	GroupID   string
}

type TaskRepository interface {
	GetByID(ctx context.Context, id string) (Task, error)
	ListByGroup(ctx context.Context, groupID string) ([]Task, error)
	ListByAssignee(ctx context.Context, userID, groupID string) ([]Task, error)
	Create(ctx context.Context, params CreateTaskParams) (Task, error)
	Update(ctx context.Context, id string, params UpdateTaskParams) (Task, error)
	Delete(ctx context.Context, id string) error
}

type TaskOccurrenceRepository interface {
	GetByID(ctx context.Context, id string) (TaskOccurrence, error)
	ListByTask(ctx context.Context, taskID string) ([]TaskOccurrence, error)
	ListPendingByUser(ctx context.Context, userID string) ([]TaskOccurrenceWithTask, error)
	ListByDateRange(ctx context.Context, groupID string, start, end time.Time) ([]TaskOccurrenceWithTask, error)
	Create(ctx context.Context, taskID string, dueDate time.Time, status string) (TaskOccurrence, error)
	Complete(ctx context.Context, id string, completedBy string) error
	Discard(ctx context.Context, id string) error
	MarkAsOverdue(ctx context.Context, id string) error
	DeleteByTask(ctx context.Context, taskID string) error
}

var (
	ErrTaskNotFound           = errors.New("task not found")
	ErrTaskOccurrenceNotFound = errors.New("task occurrence not found")
	ErrInvalidDueDate         = errors.New("due date must be in the future")
)
