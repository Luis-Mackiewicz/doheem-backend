package http

import (
	"net/http"
	"time"

	"doheem-backend/internal/domain"
)

type TaskHandler struct {
	svc *domain.TaskService
}

func NewTaskHandler(svc *domain.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	groupID := r.PathValue("groupId")

	var req struct {
		Title            string  `json:"title"              validate:"required"`
		Description      *string `json:"description,omitempty"`
		AssignedTo       *string `json:"assigned_to,omitempty"`
		Category         *string `json:"category,omitempty"`
		IsRecurring      bool    `json:"is_recurring"`
		RecurringPeriod  *string `json:"recurring_period,omitempty"`
		RecurringEndedAt *string `json:"recurring_ended_at,omitempty"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}

	var recurringEndedAt *time.Time
	if req.RecurringEndedAt != nil {
		t, err := time.Parse("2006-01-02T15:04:05Z07:00", *req.RecurringEndedAt)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid recurring_ended_at format")
			return
		}
		recurringEndedAt = &t
	}

	task, err := h.svc.Create(r.Context(), domain.CreateTaskParams{
		GroupID:          groupID,
		Title:            req.Title,
		Description:      req.Description,
		AssignedTo:       req.AssignedTo,
		Category:         req.Category,
		IsRecurring:      req.IsRecurring,
		RecurringPeriod:  req.RecurringPeriod,
		RecurringEndedAt: recurringEndedAt,
		CreatedBy:        userID,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toTaskResponse(task))
}

func (h *TaskHandler) ListByGroup(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupId")
	tasks, err := h.svc.ListByGroup(r.Context(), groupID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toTaskResponses(tasks))
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toTaskResponse(task))
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		Title            *string `json:"title,omitempty"`
		Description      *string `json:"description,omitempty"`
		AssignedTo       *string `json:"assigned_to,omitempty"`
		Category         *string `json:"category,omitempty"`
		IsRecurring      *bool   `json:"is_recurring,omitempty"`
		RecurringPeriod  *string `json:"recurring_period,omitempty"`
		RecurringEndedAt *string `json:"recurring_ended_at,omitempty"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	var recurringEndedAt *time.Time
	if req.RecurringEndedAt != nil {
		t, err := time.Parse("2006-01-02T15:04:05Z07:00", *req.RecurringEndedAt)
		if err == nil {
			recurringEndedAt = &t
		}
	}
	task, err := h.svc.Update(r.Context(), id, domain.UpdateTaskParams{
		Title:            req.Title,
		Description:      req.Description,
		AssignedTo:       req.AssignedTo,
		Category:         req.Category,
		IsRecurring:      req.IsRecurring,
		RecurringPeriod:  req.RecurringPeriod,
		RecurringEndedAt: recurringEndedAt,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toTaskResponse(task))
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) CreateOccurrence(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskId")
	var req struct {
		DueDate string `json:"due_date" validate:"required"`
		Status  string `json:"status"   validate:"required,oneof=pending completed discarded overdue"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid due_date, use YYYY-MM-DD")
		return
	}
	occurrence, err := h.svc.CreateOccurrence(r.Context(), taskID, dueDate, req.Status)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toTaskOccurrenceResponse(occurrence))
}

func (h *TaskHandler) ListOccurrences(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")
	occurrences, err := h.svc.ListOccurrencesByTask(r.Context(), taskID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toTaskOccurrenceResponses(occurrences))
}

func (h *TaskHandler) CompleteOccurrence(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	id := r.PathValue("id")
	if err := h.svc.CompleteOccurrence(r.Context(), id, userID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) DiscardOccurrence(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.DiscardOccurrence(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type taskResponse struct {
	ID               string   `json:"id"`
	GroupID          string   `json:"group_id"`
	Title            string   `json:"title"`
	Description      *string  `json:"description,omitempty"`
	AssignedTo       *string  `json:"assigned_to,omitempty"`
	Category         *string  `json:"category,omitempty"`
	IsRecurring      bool     `json:"is_recurring"`
	RecurringPeriod  *string  `json:"recurring_period,omitempty"`
	RecurringEndedAt *string  `json:"recurring_ended_at,omitempty"`
	CreatedBy        string   `json:"created_by"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
}

type taskOccurrenceResponse struct {
	ID          string  `json:"id"`
	TaskID      string  `json:"task_id"`
	DueDate     string  `json:"due_date"`
	Status      string  `json:"status"`
	CompletedBy *string `json:"completed_by,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

func toTaskResponse(t domain.Task) taskResponse {
	r := taskResponse{
		ID:          t.ID,
		GroupID:     t.GroupID,
		Title:       t.Title,
		Description: t.Description,
		AssignedTo:  t.AssignedTo,
		Category:    t.Category,
		IsRecurring: t.IsRecurring,
		CreatedBy:   t.CreatedBy,
		CreatedAt:   t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   t.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if t.RecurringPeriod != nil {
		r.RecurringPeriod = t.RecurringPeriod
	}
	if t.RecurringEndedAt != nil {
		s := t.RecurringEndedAt.Format("2006-01-02T15:04:05Z07:00")
		r.RecurringEndedAt = &s
	}
	return r
}

func toTaskResponses(tasks []domain.Task) []taskResponse {
	res := make([]taskResponse, len(tasks))
	for i, t := range tasks {
		res[i] = toTaskResponse(t)
	}
	return res
}

func toTaskOccurrenceResponse(o domain.TaskOccurrence) taskOccurrenceResponse {
	return taskOccurrenceResponse{
		ID:          o.ID,
		TaskID:      o.TaskID,
		DueDate:     o.DueDate.Format("2006-01-02"),
		Status:      o.Status,
		CompletedBy: o.CompletedBy,
		CreatedAt:   o.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toTaskOccurrenceResponses(occurrences []domain.TaskOccurrence) []taskOccurrenceResponse {
	res := make([]taskOccurrenceResponse, len(occurrences))
	for i, o := range occurrences {
		res[i] = toTaskOccurrenceResponse(o)
	}
	return res
}
