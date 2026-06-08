package task

import (
	"context"
	"time"

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
	t, err := r.q.CreateTask(ctx, db.CreateTaskParams{
		GroupID:     db.UUIDFromString(params.GroupID),
		Title:       params.Title,
		Description: params.Description,
		AssignedTo:  db.UUIDFromString(params.AssignedTo),
		CreatedBy:   db.UUIDFromString(params.CreatedBy),
		DueDate:     db.DateFromTime(params.DueDate),
	})
	if err != nil {
		return Task{}, err
	}
	return toTask(t), nil
}

func (r *TaskRepo) Update(ctx context.Context, id string, params UpdateTaskParams) (Task, error) {
	t, err := r.q.UpdateTask(ctx, db.UpdateTaskParams{
		ID:          db.UUIDFromString(id),
		Title:       deptrStr(params.Title),
		Description: deptrStr(params.Description),
		AssignedTo:  db.UUIDFromStringPtr(params.AssignedTo),
		Status:      deptrStr(params.Status),
		Position:    deptrInt32(params.Position),
		DueDate:     deptrDate(params.DueDate),
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
		ID:          db.UUIDToString(t.ID),
		GroupID:     db.UUIDToString(t.GroupID),
		Title:       t.Title,
		Description: t.Description,
		AssignedTo:  db.UUIDToString(t.AssignedTo),
		CreatedBy:   db.UUIDToString(t.CreatedBy),
		Status:      t.Status,
		Position:    t.Position,
		DueDate:     t.DueDate.Time,
		CreatedAt:   t.CreatedAt.Time,
		UpdatedAt:   t.UpdatedAt.Time,
	}
}

func toTasks(tasks []db.Task) []Task {
	result := make([]Task, len(tasks))
	for i, t := range tasks {
		result[i] = toTask(t)
	}
	return result
}

func deptrStr(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func deptrInt32(i *int32) int32 {
	if i != nil {
		return *i
	}
	return 0
}

func deptrDate(t *time.Time) pgtype.Date {
	if t != nil {
		return db.DateFromTime(*t)
	}
	return pgtype.Date{}
}
