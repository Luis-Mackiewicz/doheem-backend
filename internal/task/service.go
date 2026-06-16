package task

import (
	"context"
	"fmt"
	"time"

	"doheem-backend/internal/group"
	"doheem-backend/internal/notification"
)

var validStatusTransitions = map[string]map[string]bool{
	"todo":  {"todo": true, "doing": true},
	"doing": {"doing": true, "todo": true, "done": true},
	"done":  {"done": true, "doing": true},
}

type TaskService struct {
	taskRepo         TaskRepository
	taskOccurrenceRepo TaskOccurrenceRepository
	memberRepo       group.GroupMemberRepository
	notifRepo        notification.NotificationRepository
}

func NewTaskService(
	taskRepo TaskRepository,
	taskOccurrenceRepo TaskOccurrenceRepository,
	memberRepo group.GroupMemberRepository,
	notifRepo notification.NotificationRepository,
) *TaskService {
	return &TaskService{
		taskRepo:         taskRepo,
		taskOccurrenceRepo: taskOccurrenceRepo,
		memberRepo:       memberRepo,
		notifRepo:        notifRepo,
	}
}

func (s *TaskService) Create(ctx context.Context, params CreateTaskParams) (Task, error) {
	if params.DueDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return Task{}, ErrInvalidDueDate
	}

	task, err := s.taskRepo.Create(ctx, params)
	if err != nil {
		return Task{}, err
	}

	if params.AssignedTo != params.CreatedBy {
		title := fmt.Sprintf("Nova tarefa: %s", params.Title)
		message := fmt.Sprintf("Você foi atribuído à tarefa \"%s\" — vence em %s", params.Title, params.DueDate.Format("02/01/2006"))
		relatedID := &task.ID
		s.notifRepo.Create(ctx, notification.CreateNotificationParams{
			UserID:    params.AssignedTo,
			Type:      "task",
			Title:     title,
			Message:   message,
			RelatedID: relatedID,
		})
	}

	return task, nil
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

func (s *TaskService) Update(ctx context.Context, id string, params UpdateTaskParams, userID string) (Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return Task{}, ErrTaskNotFound
	}

	if task.CreatedBy != userID {
		_, err := s.memberRepo.Get(ctx, task.GroupID, userID)
		if err != nil || !s.isAdminOrCreator(task, userID) {
			return Task{}, ErrForbidden
		}
	}

	if params.Status != nil {
		current := task.Status
		next := *params.Status
		if !validStatusTransitions[current][next] {
			return Task{}, ErrInvalidStatusTransition
		}
	}

	return s.taskRepo.Update(ctx, id, params)
}

func (s *TaskService) Delete(ctx context.Context, id string, userID string) error {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return ErrTaskNotFound
	}

	if !s.isAdminOrCreator(task, userID) && task.CreatedBy != userID {
		return ErrForbidden
	}

	return s.taskRepo.Delete(ctx, id)
}

func (s *TaskService) isAdminOrCreator(task Task, userID string) bool {
	if task.CreatedBy == userID {
		return true
	}
	member, err := s.memberRepo.Get(context.Background(), task.GroupID, userID)
	if err != nil {
		return false
	}
	return member.IsAdmin
}

func (s *TaskService) CanModify(task Task, userID string) bool {
	return s.isAdminOrCreator(task, userID)
}

func (s *TaskService) IsAdmin(ctx context.Context, groupID, userID string) bool {
	member, err := s.memberRepo.Get(ctx, groupID, userID)
	if err != nil {
		return false
	}
	return member.IsAdmin
}

func (s *TaskService) CreateOccurrence(ctx context.Context, taskID string, dueDate time.Time, status string) (TaskOccurrence, error) {
	return s.taskOccurrenceRepo.Create(ctx, taskID, dueDate, status)
}

func (s *TaskService) CompleteOccurrence(ctx context.Context, id, completedBy string) error {
	if err := s.taskOccurrenceRepo.Complete(ctx, id, completedBy); err != nil {
		return err
	}
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
