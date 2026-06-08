package task

import (
	"context"

	"doheem-backend/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
)

type TaskRepo struct {
	q *db.Queries
}

func NewTaskRepo(q *db.Queries) *TaskRepo {
	return &TaskRepo{q: q}
}

func (r *TaskRepo) GetByID(ctx context.Context, id string) (Task, error) {
	t, err := r.q.GetTaskByID(ctx, db.UUIDFromString(id))
	if err != nil {
		return Task{}, err
	}
	return toTask(t), nil
}

func (r *TaskRepo) ListByGroup(ctx context.Context, groupID string) ([]Task, error) {
	tasks, err := r.q.ListTasksByGroup(ctx, db.UUIDFromString(groupID))
	if err != nil {
		return nil, err
	}
	return toTasks(tasks), nil
}

func (r *TaskRepo) ListByAssignee(ctx context.Context, userID, groupID string) ([]Task, error) {
	tasks, err := r.q.ListTasksByAssignee(ctx, db.ListTasksByAssigneeParams{
		AssignedTo: db.UUIDFromString(userID),
		GroupID:    db.UUIDFromString(groupID),
	})
	if err != nil {
		return nil, err
	}
	return toTasks(tasks), nil
}

func (r *TaskRepo) Create(ctx context.Context, params CreateTaskParams) (Task, error) {
	var recurringEndedAt pgtype.Timestamptz
	if params.RecurringEndedAt != nil {
		recurringEndedAt = db.TimestamptzFromTime(*params.RecurringEndedAt)
	}
	t, err := r.q.CreateTask(ctx, db.CreateTaskParams{
		GroupID:          db.UUIDFromString(params.GroupID),
		Title:            params.Title,
		Description:      db.TextFromStringPtr(params.Description),
		AssignedTo:       db.UUIDFromStringPtr(params.AssignedTo),
		Category:         db.TextFromStringPtr(params.Category),
		IsRecurring:      params.IsRecurring,
		RecurringPeriod:  db.TextFromStringPtr(params.RecurringPeriod),
		RecurringEndedAt: recurringEndedAt,
		CreatedBy:        db.UUIDFromString(params.CreatedBy),
	})
	if err != nil {
		return Task{}, err
	}
	return toTask(t), nil
}

func (r *TaskRepo) Update(ctx context.Context, id string, params UpdateTaskParams) (Task, error) {
	var title string
	if params.Title != nil {
		title = *params.Title
	}
	var isRecurring bool
	if params.IsRecurring != nil {
		isRecurring = *params.IsRecurring
	}
	var recurringEndedAt pgtype.Timestamptz
	if params.RecurringEndedAt != nil {
		recurringEndedAt = db.TimestamptzFromTime(*params.RecurringEndedAt)
	}
	t, err := r.q.UpdateTask(ctx, db.UpdateTaskParams{
		ID:               db.UUIDFromString(id),
		Title:            title,
		Description:      db.TextFromStringPtr(params.Description),
		AssignedTo:       db.UUIDFromStringPtr(params.AssignedTo),
		Category:         db.TextFromStringPtr(params.Category),
		IsRecurring:      isRecurring,
		RecurringPeriod:  db.TextFromStringPtr(params.RecurringPeriod),
		RecurringEndedAt: recurringEndedAt,
	})
	if err != nil {
		return Task{}, err
	}
	return toTask(t), nil
}

func (r *TaskRepo) Delete(ctx context.Context, id string) error {
	return r.q.DeleteTask(ctx, db.UUIDFromString(id))
}

func toTask(t db.Task) Task {
	return Task{
		ID:               db.UUIDToString(t.ID),
		GroupID:          db.UUIDToString(t.GroupID),
		Title:            t.Title,
		Description:      db.TextToStringPtr(t.Description),
		AssignedTo:       db.UUIDToStringPtr(t.AssignedTo),
		Category:         db.TextToStringPtr(t.Category),
		IsRecurring:      t.IsRecurring,
		RecurringPeriod:  db.TextToStringPtr(t.RecurringPeriod),
		RecurringEndedAt: db.TimestamptzToTimePtr(t.RecurringEndedAt),
		CreatedBy:        db.UUIDToString(t.CreatedBy),
		CreatedAt:        t.CreatedAt.Time,
		UpdatedAt:        t.UpdatedAt.Time,
	}
}

func toTasks(tasks []db.Task) []Task {
	result := make([]Task, len(tasks))
	for i, t := range tasks {
		result[i] = toTask(t)
	}
	return result
}
