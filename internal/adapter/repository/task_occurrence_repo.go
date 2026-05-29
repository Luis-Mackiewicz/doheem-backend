package repository

import (
	"context"
	"time"

	"doheem-backend/internal/adapter/repository/db"
	"doheem-backend/internal/domain"
)

type TaskOccurrenceRepo struct {
	q *db.Queries
}

func NewTaskOccurrenceRepo(q *db.Queries) *TaskOccurrenceRepo {
	return &TaskOccurrenceRepo{q: q}
}

func (r *TaskOccurrenceRepo) GetByID(ctx context.Context, id string) (domain.TaskOccurrence, error) {
	to, err := r.q.GetTaskOccurrenceByID(ctx, uuidFromString(id))
	if err != nil {
		return domain.TaskOccurrence{}, err
	}
	return domainTaskOccurrence(to), nil
}

func (r *TaskOccurrenceRepo) ListByTask(ctx context.Context, taskID string) ([]domain.TaskOccurrence, error) {
	occurrences, err := r.q.ListTaskOccurrencesByTask(ctx, uuidFromString(taskID))
	if err != nil {
		return nil, err
	}
	result := make([]domain.TaskOccurrence, len(occurrences))
	for i, o := range occurrences {
		result[i] = domainTaskOccurrence(o)
	}
	return result, nil
}

func (r *TaskOccurrenceRepo) ListPendingByUser(ctx context.Context, userID string) ([]domain.TaskOccurrenceWithTask, error) {
	rows, err := r.q.ListPendingTaskOccurrencesByUser(ctx, uuidFromString(userID))
	if err != nil {
		return nil, err
	}
	return domainTaskOccurrencesWithTask(rows), nil
}

func (r *TaskOccurrenceRepo) ListByDateRange(ctx context.Context, groupID string, start, end time.Time) ([]domain.TaskOccurrenceWithTask, error) {
	rows, err := r.q.ListTaskOccurrencesByDateRange(ctx, db.ListTaskOccurrencesByDateRangeParams{
		GroupID:   uuidFromString(groupID),
		DueDate:   dateFromTime(start),
		DueDate_2: dateFromTime(end),
	})
	if err != nil {
		return nil, err
	}
	result := make([]domain.TaskOccurrenceWithTask, len(rows))
	for i, row := range rows {
		result[i] = domain.TaskOccurrenceWithTask{
			TaskOccurrence: domain.TaskOccurrence{
				ID:          uuidToString(row.ID),
				TaskID:      uuidToString(row.TaskID),
				DueDate:     row.DueDate.Time,
				Status:      row.Status,
				CompletedBy: uuidToStringPtr(row.CompletedBy),
				CompletedAt: timestamptzToTimePtr(row.CompletedAt),
				DiscardedAt: timestamptzToTimePtr(row.DiscardedAt),
				CreatedAt:   row.CreatedAt.Time,
			},
			TaskTitle: row.TaskTitle,
			GroupID:   uuidToString(row.GroupID),
		}
	}
	return result, nil
}

func (r *TaskOccurrenceRepo) Create(ctx context.Context, taskID string, dueDate time.Time, status string) (domain.TaskOccurrence, error) {
	to, err := r.q.CreateTaskOccurrence(ctx, db.CreateTaskOccurrenceParams{
		TaskID:  uuidFromString(taskID),
		DueDate: dateFromTime(dueDate),
		Status:  status,
	})
	if err != nil {
		return domain.TaskOccurrence{}, err
	}
	return domainTaskOccurrence(to), nil
}

func (r *TaskOccurrenceRepo) Complete(ctx context.Context, id string, completedBy string) error {
	return r.q.CompleteTaskOccurrence(ctx, db.CompleteTaskOccurrenceParams{
		ID:          uuidFromString(id),
		CompletedBy: uuidFromString(completedBy),
	})
}

func (r *TaskOccurrenceRepo) Discard(ctx context.Context, id string) error {
	return r.q.DiscardTaskOccurrence(ctx, uuidFromString(id))
}

func (r *TaskOccurrenceRepo) MarkAsOverdue(ctx context.Context, id string) error {
	return r.q.MarkTaskOccurrenceAsOverdue(ctx, uuidFromString(id))
}

func (r *TaskOccurrenceRepo) DeleteByTask(ctx context.Context, taskID string) error {
	return r.q.DeleteTaskOccurrencesByTask(ctx, uuidFromString(taskID))
}
