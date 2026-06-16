package http

import (
	"net/http"

	"doheem-backend/internal/notification"
)

type NotificationHandler struct {
	svc *notification.NotificationService
}

func NewNotificationHandler(svc *notification.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	var req struct {
		Type    string  `json:"type"    validate:"required"`
		Title   string  `json:"title"   validate:"required"`
		Message string  `json:"message" validate:"required"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	n, err := h.svc.Create(r.Context(), notification.CreateNotificationParams{
		UserID:  userID,
		Type:    req.Type,
		Title:   req.Title,
		Message: req.Message,
	})
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusCreated, toNotificationResponse(n))
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	limit, offset := parsePagination(r)
	search := r.URL.Query().Get("search")
	notifications, total, err := h.svc.ListByUser(r.Context(), userID, search, int32(limit), int32(offset))
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, paginatedResponse{Data: toNotificationResponses(notifications), Total: int(total)})
}

func (h *NotificationHandler) ListUnread(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	notifications, err := h.svc.ListUnreadByUser(r.Context(), userID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, toNotificationResponses(notifications))
}

func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	id := r.PathValue("id")
	if err := h.svc.MarkAsRead(r.Context(), id, userID); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *NotificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	if err := h.svc.MarkAllAsRead(r.Context(), userID); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *NotificationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	id := r.PathValue("id")
	if err := h.svc.Delete(r.Context(), id, userID); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *NotificationHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	if err := h.svc.DeleteAll(r.Context(), userID); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type notificationResponse struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Type      string  `json:"type"`
	Title     string  `json:"title"`
	Message   string  `json:"message"`
	IsRead    bool    `json:"is_read"`
	RelatedID *string `json:"related_id,omitempty"`
	CreatedAt string  `json:"created_at"`
}

func toNotificationResponse(n notification.Notification) notificationResponse {
	return notificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		Type:      n.Type,
		Title:     n.Title,
		Message:   n.Message,
		IsRead:    n.IsRead,
		RelatedID: n.RelatedID,
		CreatedAt: n.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toNotificationResponses(notifications []notification.Notification) []notificationResponse {
	res := make([]notificationResponse, len(notifications))
	for i, n := range notifications {
		res[i] = toNotificationResponse(n)
	}
	return res
}
