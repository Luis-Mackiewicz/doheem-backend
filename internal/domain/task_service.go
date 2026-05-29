package domain

import (
	"context"
	"time"
)

type TaskService struct {
	taskRepo         TaskRepository
	taskOccurrenceRepo TaskOccurrenceRepository
	eventBus         EventBus
}

func NewTaskService(taskRepo TaskRepository, taskOccurrenceRepo TaskOccurrenceRepository, eventBus EventBus) *TaskService {
	return &TaskService{
		taskRepo:         taskRepo,
		taskOccurrenceRepo: taskOccurrenceRepo,
		eventBus:         eventBus,
	}
}

func (s *TaskService) Create(ctx context.Context, params CreateTaskParams) (Task, error) {
	return s.taskRepo.Create(ctx, params)
}

func (s *TaskService) GetByID(ctx context.Context, id string) (Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return Task{}, ErrTaskNotFound
	}
	return task, nil
}

func (s *TaskService) ListByGroup(ctx context.Context, groupID string) ([]Task, error) {
	return s.taskRepo.ListByGroup(ctx, groupID)
}

func (s *TaskService) ListByAssignee(ctx context.Context, userID, groupID string) ([]Task, error) {
	return s.taskRepo.ListByAssignee(ctx, userID, groupID)
}

func (s *TaskService) Update(ctx context.Context, id string, params UpdateTaskParams) (Task, error) {
	return s.taskRepo.Update(ctx, id, params)
}

func (s *TaskService) Delete(ctx context.Context, id string) error {
	return s.taskRepo.Delete(ctx, id)
}

func (s *TaskService) CreateOccurrence(ctx context.Context, taskID string, dueDate time.Time, status string) (TaskOccurrence, error) {
	return s.taskOccurrenceRepo.Create(ctx, taskID, dueDate, status)
}

func (s *TaskService) CompleteOccurrence(ctx context.Context, id, completedBy string) error {
	if err := s.taskOccurrenceRepo.Complete(ctx, id, completedBy); err != nil {
		return err
	}

	occ, err := s.taskOccurrenceRepo.GetByID(ctx, id)
	if err != nil {
		return nil
	}
	task, err := s.taskRepo.GetByID(ctx, occ.TaskID)
	if err != nil {
		return nil
	}

	s.eventBus.Publish(ctx, DomainEvent{
		Type:     "task.occurrence.completed",
		EntityID: id,
		UserID:   completedBy,
		GroupID:  task.GroupID,
		Payload: map[string]any{
			"task_id":    task.ID,
			"task_title": task.Title,
			"assigned_to": task.AssignedTo,
		},
	})
	return nil
}

func (s *TaskService) DiscardOccurrence(ctx context.Context, id string) error {
	return s.taskOccurrenceRepo.Discard(ctx, id)
}

func (s *TaskService) MarkOccurrenceAsOverdue(ctx context.Context, id string) error {
	return s.taskOccurrenceRepo.MarkAsOverdue(ctx, id)
}

func (s *TaskService) ListPendingOccurrences(ctx context.Context, userID string) ([]TaskOccurrenceWithTask, error) {
	return s.taskOccurrenceRepo.ListPendingByUser(ctx, userID)
}

func (s *TaskService) ListOccurrencesByDateRange(ctx context.Context, groupID string, start, end time.Time) ([]TaskOccurrenceWithTask, error) {
	return s.taskOccurrenceRepo.ListByDateRange(ctx, groupID, start, end)
}

func (s *TaskService) ListOccurrencesByTask(ctx context.Context, taskID string) ([]TaskOccurrence, error) {
	return s.taskOccurrenceRepo.ListByTask(ctx, taskID)
}
