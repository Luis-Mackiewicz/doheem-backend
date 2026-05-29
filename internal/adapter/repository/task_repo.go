package repository

import (
	"context"

	"doheem-backend/internal/adapter/repository/db"
	"doheem-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgtype"
)

type TaskRepo struct {
	q *db.Queries
}

func NewTaskRepo(q *db.Queries) *TaskRepo {
	return &TaskRepo{q: q}
}

func (r *TaskRepo) GetByID(ctx context.Context, id string) (domain.Task, error) {
	t, err := r.q.GetTaskByID(ctx, uuidFromString(id))
	if err != nil {
		return domain.Task{}, err
	}
	return domainTask(t), nil
}

func (r *TaskRepo) ListByGroup(ctx context.Context, groupID string) ([]domain.Task, error) {
	tasks, err := r.q.ListTasksByGroup(ctx, uuidFromString(groupID))
	if err != nil {
		return nil, err
	}
	return domainTasks(tasks), nil
}

func (r *TaskRepo) ListByAssignee(ctx context.Context, userID, groupID string) ([]domain.Task, error) {
	tasks, err := r.q.ListTasksByAssignee(ctx, db.ListTasksByAssigneeParams{
		AssignedTo: uuidFromString(userID),
		GroupID:    uuidFromString(groupID),
	})
	if err != nil {
		return nil, err
	}
	return domainTasks(tasks), nil
}

func (r *TaskRepo) Create(ctx context.Context, params domain.CreateTaskParams) (domain.Task, error) {
	var recurringEndedAt pgtype.Timestamptz
	if params.RecurringEndedAt != nil {
		recurringEndedAt = timestamptzFromTime(*params.RecurringEndedAt)
	}
	t, err := r.q.CreateTask(ctx, db.CreateTaskParams{
		GroupID:          uuidFromString(params.GroupID),
		Title:            params.Title,
		Description:      textFromStringPtr(params.Description),
		AssignedTo:       uuidFromStringPtr(params.AssignedTo),
		Category:         textFromStringPtr(params.Category),
		IsRecurring:      params.IsRecurring,
		RecurringPeriod:  textFromStringPtr(params.RecurringPeriod),
		RecurringEndedAt: recurringEndedAt,
		CreatedBy:        uuidFromString(params.CreatedBy),
	})
	if err != nil {
		return domain.Task{}, err
	}
	return domainTask(t), nil
}

func (r *TaskRepo) Update(ctx context.Context, id string, params domain.UpdateTaskParams) (domain.Task, error) {
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
		recurringEndedAt = timestamptzFromTime(*params.RecurringEndedAt)
	}
	t, err := r.q.UpdateTask(ctx, db.UpdateTaskParams{
		ID:               uuidFromString(id),
		Title:            title,
		Description:      textFromStringPtr(params.Description),
		AssignedTo:       uuidFromStringPtr(params.AssignedTo),
		Category:         textFromStringPtr(params.Category),
		IsRecurring:      isRecurring,
		RecurringPeriod:  textFromStringPtr(params.RecurringPeriod),
		RecurringEndedAt: recurringEndedAt,
	})
	if err != nil {
		return domain.Task{}, err
	}
	return domainTask(t), nil
}

func (r *TaskRepo) Delete(ctx context.Context, id string) error {
	return r.q.DeleteTask(ctx, uuidFromString(id))
}
