package http

import (
	"net/http"
	"time"

	"doheem-backend/internal/task"
)

type TaskHandler struct {
	svc *task.TaskService
}

func NewTaskHandler(svc *task.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	groupID := r.PathValue("groupId")

	var req struct {
		Title       string `json:"title"        validate:"required"`
		Description string `json:"description,omitempty"`
		AssignedTo  string `json:"assigned_to"   validate:"required"`
		DueDate     string `json:"due_date"      validate:"required"`
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

	created, err := h.svc.Create(r.Context(), task.CreateTaskParams{
		GroupID:     groupID,
		Title:       req.Title,
		Description: req.Description,
		AssignedTo:  req.AssignedTo,
		CreatedBy:   userID,
		DueDate:     dueDate,
	})
	if err != nil {
		handleError(w, r, err)
		return
	}
	isAdmin := h.svc.IsAdmin(r.Context(), groupID, userID)
	respondJSON(w, http.StatusCreated, toTaskResponse(created, userID, isAdmin))
}

func (h *TaskHandler) ListByGroup(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	groupID := r.PathValue("groupId")
	tasks, err := h.svc.ListByGroup(r.Context(), groupID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	isAdmin := h.svc.IsAdmin(r.Context(), groupID, userID)
	limit, offset := parsePagination(r)
	items, total := paginate(toTaskResponses(tasks, userID, isAdmin), limit, offset)
	respondJSON(w, http.StatusOK, paginatedResponse{Data: items, Total: total})
}

func (h *TaskHandler) ListOccurrences(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")
	occurrences, err := h.svc.ListOccurrencesByTask(r.Context(), taskID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	limit, offset := parsePagination(r)
	items, total := paginate(toTaskOccurrenceResponses(occurrences), limit, offset)
	respondJSON(w, http.StatusOK, paginatedResponse{Data: items, Total: total})
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	id := r.PathValue("id")
	t, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, r, err)
		return
	}
	isAdmin := h.svc.IsAdmin(r.Context(), t.GroupID, userID)
	respondJSON(w, http.StatusOK, toTaskResponse(t, userID, isAdmin))
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	id := r.PathValue("id")
	var req struct {
		Title       *string `json:"title,omitempty"`
		Description *string `json:"description,omitempty"`
		AssignedTo  *string `json:"assigned_to,omitempty"`
		Status      *string `json:"status,omitempty"`
		Position    *int32  `json:"position,omitempty"`
		DueDate     *string `json:"due_date,omitempty"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}

	params := task.UpdateTaskParams{
		Title:       req.Title,
		Description: req.Description,
		AssignedTo:  req.AssignedTo,
		Status:      req.Status,
		Position:    req.Position,
	}
	if req.DueDate != nil {
		t, err := time.Parse("2006-01-02", *req.DueDate)
		if err == nil {
			params.DueDate = &t
		}
	}

	updated, err := h.svc.Update(r.Context(), id, params, userID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	isAdmin := h.svc.IsAdmin(r.Context(), updated.GroupID, userID)
	respondJSON(w, http.StatusOK, toTaskResponse(updated, userID, isAdmin))
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	id := r.PathValue("id")
	if err := h.svc.Delete(r.Context(), id, userID); err != nil {
		handleError(w, r, err)
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
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusCreated, toTaskOccurrenceResponse(occurrence))
}

func (h *TaskHandler) CompleteOccurrence(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	id := r.PathValue("id")
	if err := h.svc.CompleteOccurrence(r.Context(), id, userID); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) DiscardOccurrence(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.DiscardOccurrence(r.Context(), id); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type taskResponse struct {
	ID          string `json:"id"`
	GroupID     string `json:"group_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	AssignedTo  string `json:"assigned_to"`
	CreatedBy   string `json:"created_by"`
	Status      string `json:"status"`
	Position    int32  `json:"position"`
	DueDate     string `json:"due_date"`
	IsOverdue   bool   `json:"is_overdue"`
	CanModify   bool   `json:"can_modify"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type taskOccurrenceResponse struct {
	ID          string  `json:"id"`
	TaskID      string  `json:"task_id"`
	DueDate     string  `json:"due_date"`
	Status      string  `json:"status"`
	CompletedBy *string `json:"completed_by,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

func toTaskResponse(t task.Task, currentUser string, isAdmin bool) taskResponse {
	now := time.Now().Truncate(24 * time.Hour)
	isOverdue := t.Status != "done" && !t.DueDate.IsZero() && t.DueDate.Before(now)
	canModify := t.CreatedBy == currentUser || isAdmin
	return taskResponse{
		ID:          t.ID,
		GroupID:     t.GroupID,
		Title:       t.Title,
		Description: t.Description,
		AssignedTo:  t.AssignedTo,
		CreatedBy:   t.CreatedBy,
		Status:      t.Status,
		Position:    t.Position,
		DueDate:     t.DueDate.Format("2006-01-02"),
		IsOverdue:   isOverdue,
		CanModify:   canModify,
		CreatedAt:   t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   t.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toTaskResponses(tasks []task.Task, currentUser string, isAdmin bool) []taskResponse {
	res := make([]taskResponse, len(tasks))
	for i, t := range tasks {
		res[i] = toTaskResponse(t, currentUser, isAdmin)
	}
	return res
}

func toTaskOccurrenceResponse(o task.TaskOccurrence) taskOccurrenceResponse {
	return taskOccurrenceResponse{
		ID:          o.ID,
		TaskID:      o.TaskID,
		DueDate:     o.DueDate.Format("2006-01-02"),
		Status:      o.Status,
		CompletedBy: o.CompletedBy,
		CreatedAt:   o.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toTaskOccurrenceResponses(occurrences []task.TaskOccurrence) []taskOccurrenceResponse {
	res := make([]taskOccurrenceResponse, len(occurrences))
	for i, o := range occurrences {
		res[i] = toTaskOccurrenceResponse(o)
	}
	return res
}
