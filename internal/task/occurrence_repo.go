package task

import (
	"context"
	"time"

	"doheem-backend/internal/db"
)

type TaskOccurrenceRepo struct {
	q *db.Queries
}

func NewTaskOccurrenceRepo(q *db.Queries) *TaskOccurrenceRepo {
	return &TaskOccurrenceRepo{q: q}
}

func (r *TaskOccurrenceRepo) GetByID(ctx context.Context, id string) (TaskOccurrence, error) {
	to, err := r.q.GetTaskOccurrenceByID(ctx, db.UUIDFromString(id))
	if err != nil {
		return TaskOccurrence{}, err
	}
	return toOccurrence(to), nil
}

func (r *TaskOccurrenceRepo) ListByTask(ctx context.Context, taskID string) ([]TaskOccurrence, error) {
	occurrences, err := r.q.ListTaskOccurrencesByTask(ctx, db.UUIDFromString(taskID))
	if err != nil {
		return nil, err
	}
	result := make([]TaskOccurrence, len(occurrences))
	for i, o := range occurrences {
		result[i] = toOccurrence(o)
	}
	return result, nil
}

func (r *TaskOccurrenceRepo) ListPendingByUser(ctx context.Context, userID string) ([]TaskOccurrenceWithTask, error) {
	rows, err := r.q.ListPendingTaskOccurrencesByUser(ctx, db.UUIDFromString(userID))
	if err != nil {
		return nil, err
	}
	return toOccurrencesWithTask(rows), nil
}

func (r *TaskOccurrenceRepo) ListByDateRange(ctx context.Context, groupID string, start, end time.Time) ([]TaskOccurrenceWithTask, error) {
	rows, err := r.q.ListTaskOccurrencesByDateRange(ctx, db.ListTaskOccurrencesByDateRangeParams{
		GroupID:   db.UUIDFromString(groupID),
		DueDate:   db.DateFromTime(start),
		DueDate_2: db.DateFromTime(end),
	})
	if err != nil {
		return nil, err
	}
	result := make([]TaskOccurrenceWithTask, len(rows))
	for i, row := range rows {
		result[i] = TaskOccurrenceWithTask{
			TaskOccurrence: TaskOccurrence{
				ID:          db.UUIDToString(row.ID),
				TaskID:      db.UUIDToString(row.TaskID),
				DueDate:     row.DueDate.Time,
				Status:      row.Status,
				CompletedBy: db.UUIDToStringPtr(row.CompletedBy),
				CompletedAt: db.TimestamptzToTimePtr(row.CompletedAt),
				DiscardedAt: db.TimestamptzToTimePtr(row.DiscardedAt),
				CreatedAt:   row.CreatedAt.Time,
			},
			TaskTitle: row.TaskTitle,
			GroupID:   db.UUIDToString(row.GroupID),
		}
	}
	return result, nil
}

func (r *TaskOccurrenceRepo) Create(ctx context.Context, taskID string, dueDate time.Time, status string) (TaskOccurrence, error) {
	to, err := r.q.CreateTaskOccurrence(ctx, db.CreateTaskOccurrenceParams{
		TaskID:  db.UUIDFromString(taskID),
		DueDate: db.DateFromTime(dueDate),
		Status:  status,
	})
	if err != nil {
		return TaskOccurrence{}, err
	}
	return toOccurrence(to), nil
}

func (r *TaskOccurrenceRepo) Complete(ctx context.Context, id string, completedBy string) error {
	return r.q.CompleteTaskOccurrence(ctx, db.CompleteTaskOccurrenceParams{
		ID:          db.UUIDFromString(id),
		CompletedBy: db.UUIDFromString(completedBy),
	})
}

func (r *TaskOccurrenceRepo) Discard(ctx context.Context, id string) error {
	return r.q.DiscardTaskOccurrence(ctx, db.UUIDFromString(id))
}

func (r *TaskOccurrenceRepo) MarkAsOverdue(ctx context.Context, id string) error {
	return r.q.MarkTaskOccurrenceAsOverdue(ctx, db.UUIDFromString(id))
}

func (r *TaskOccurrenceRepo) DeleteByTask(ctx context.Context, taskID string) error {
	return r.q.DeleteTaskOccurrencesByTask(ctx, db.UUIDFromString(taskID))
}

func toOccurrence(to db.TaskOccurrence) TaskOccurrence {
	return TaskOccurrence{
		ID:          db.UUIDToString(to.ID),
		TaskID:      db.UUIDToString(to.TaskID),
		DueDate:     to.DueDate.Time,
		Status:      to.Status,
		CompletedBy: db.UUIDToStringPtr(to.CompletedBy),
		CompletedAt: db.TimestamptzToTimePtr(to.CompletedAt),
		DiscardedAt: db.TimestamptzToTimePtr(to.DiscardedAt),
		CreatedAt:   to.CreatedAt.Time,
	}
}

func toOccurrenceWithTask(row db.ListPendingTaskOccurrencesByUserRow) TaskOccurrenceWithTask {
	return TaskOccurrenceWithTask{
		TaskOccurrence: TaskOccurrence{
			ID:          db.UUIDToString(row.ID),
			TaskID:      db.UUIDToString(row.TaskID),
			DueDate:     row.DueDate.Time,
			Status:      row.Status,
			CompletedBy: db.UUIDToStringPtr(row.CompletedBy),
			CompletedAt: db.TimestamptzToTimePtr(row.CompletedAt),
			DiscardedAt: db.TimestamptzToTimePtr(row.DiscardedAt),
			CreatedAt:   row.CreatedAt.Time,
		},
		TaskTitle: row.TaskTitle,
		GroupID:   db.UUIDToString(row.GroupID),
	}
}

func toOccurrencesWithTask(rows []db.ListPendingTaskOccurrencesByUserRow) []TaskOccurrenceWithTask {
	result := make([]TaskOccurrenceWithTask, len(rows))
	for i, r := range rows {
		result[i] = toOccurrenceWithTask(r)
	}
	return result
}
