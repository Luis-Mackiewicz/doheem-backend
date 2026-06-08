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

// List lists notifications for the authenticated user
// @Summary List notifications
// @Description List all notifications for the currently authenticated user with pagination
// @Tags Notifications
// @Accept json
// @Produce json
// @Param limit query int false "Page limit"
// @Param offset query int false "Page offset"
// @Success 200 {array} notificationResponse
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/notifications [get]
// @Security BearerAuth
func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	limit, offset := parsePagination(r)
	notifications, err := h.svc.ListByUser(r.Context(), userID, int32(limit), int32(offset))
	if err != nil {
		handleError(w, r, err)
		return
	}
	items, total := paginate(toNotificationResponses(notifications), limit, offset)
	respondJSON(w, http.StatusOK, paginatedResponse{Data: items, Total: total})
}

func (h *NotificationHandler) ListUnread(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	notifications, err := h.svc.ListUnreadByUser(r.Context(), userID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	limit, offset := parsePagination(r)
	items, total := paginate(toNotificationResponses(notifications), limit, offset)
	respondJSON(w, http.StatusOK, paginatedResponse{Data: items, Total: total})
}

// MarkAsRead marks a notification as read
// @Summary Mark notification as read
// @Description Mark a single notification as read
// @Tags Notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 204 {object} nil
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Notification not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/notifications/{id}/read [patch]
// @Security BearerAuth
func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	id := r.PathValue("id")
	if err := h.svc.MarkAsRead(r.Context(), id, userID); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// MarkAllAsRead marks all notifications as read
// @Summary Mark all notifications as read
// @Description Mark all unread notifications as read for the authenticated user
// @Tags Notifications
// @Accept json
// @Produce json
// @Success 204 {object} nil
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/notifications/read-all [patch]
// @Security BearerAuth
func (h *NotificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	if err := h.svc.MarkAllAsRead(r.Context(), userID); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Delete deletes a notification
// @Summary Delete a notification
// @Description Permanently delete a notification
// @Tags Notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 204 {object} nil
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Notification not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/notifications/{id} [delete]
// @Security BearerAuth
func (h *NotificationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	id := r.PathValue("id")
	if err := h.svc.Delete(r.Context(), id, userID); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type notificationResponse struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	GroupID   *string `json:"group_id,omitempty"`
	Type      string  `json:"type"`
	Title     string  `json:"title"`
	Message   string  `json:"message"`
	IsRead    bool    `json:"is_read"`
	CreatedAt string  `json:"created_at"`
}

func toNotificationResponse(n notification.Notification) notificationResponse {
	return notificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		GroupID:   n.GroupID,
		Type:      n.Type,
		Title:     n.Title,
		Message:   n.Message,
		IsRead:    n.IsRead,
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
