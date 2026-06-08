package http

import (
	"net/http"
	"time"

	"doheem-backend/internal/invite"
)

type InviteHandler struct {
	svc *invite.InviteService
}

func NewInviteHandler(svc *invite.InviteService) *InviteHandler {
	return &InviteHandler{svc: svc}
}

// Create creates a new invite for a group
// @Summary Create an invite
// @Description Create a new invitation link/code for a group
// @Tags Invites
// @Accept json
// @Produce json
// @Param groupId path string true "Group ID"
// @Success 201 {object} inviteResponse
// @Failure 400 {object} map[string]any "Validation error"
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Group not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/groups/{groupId}/invites [post]
// @Security BearerAuth
func (h *InviteHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	groupID := r.PathValue("groupId")
	var req struct {
		ExpiresAt string `json:"expires_at" validate:"required"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	expiresAt, err := time.Parse("2006-01-02T15:04:05Z07:00", req.ExpiresAt)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid expires_at, use RFC3339")
		return
	}
	invite, err := h.svc.Create(r.Context(), groupID, userID, expiresAt)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusCreated, toInviteResponse(invite))
}

// ListByGroup lists invites for a group
// @Summary List invites by group
// @Description List all invites for a specific group
// @Tags Invites
// @Accept json
// @Produce json
// @Param groupId path string true "Group ID"
// @Success 200 {array} inviteWithCreatorResponse
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Group not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/groups/{groupId}/invites [get]
// @Security BearerAuth
func (h *InviteHandler) ListByGroup(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupId")
	invites, err := h.svc.ListByGroup(r.Context(), groupID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	limit, offset := parsePagination(r)
	items, total := paginate(toInviteWithCreatorResponses(invites), limit, offset)
	respondJSON(w, http.StatusOK, paginatedResponse{Data: items, Total: total})
}

func (h *InviteHandler) ListPending(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	invites, err := h.svc.ListPendingByUser(r.Context(), userID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	limit, offset := parsePagination(r)
	items, total := paginate(toInviteWithGroupResponses(invites), limit, offset)
	respondJSON(w, http.StatusOK, paginatedResponse{Data: items, Total: total})
}

// ListPending lists pending invites for the authenticated user
// @Summary List pending invites
// @Description List all pending invites for the currently authenticated user
// @Tags Invites
// @Accept json
// @Produce json
// @Success 200 {array} inviteWithGroupResponse
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/invites/pending [get]
// @Security BearerAuth
// Accept accepts an invite
// @Summary Accept an invite
// @Description Accept an invitation to join a group
// @Tags Invites
// @Accept json
// @Produce json
// @Param id path string true "Invite ID"
// @Success 204 {object} nil
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Invite not found"
// @Failure 409 {object} map[string]any "Invite already used or expired"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/invites/{id}/accept [post]
// @Security BearerAuth
func (h *InviteHandler) Accept(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	id := r.PathValue("id")
	if err := h.svc.Use(r.Context(), id, userID); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Revoke revokes an invite
// @Summary Revoke an invite
// @Description Revoke an invitation before it is used
// @Tags Invites
// @Accept json
// @Produce json
// @Param id path string true "Invite ID"
// @Success 204 {object} nil
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Invite not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/invites/{id}/revoke [patch]
// @Security BearerAuth
func (h *InviteHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Revoke(r.Context(), id); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type inviteResponse struct {
	ID        string  `json:"id"`
	GroupID   string  `json:"group_id"`
	Code      string  `json:"code"`
	CreatedBy string  `json:"created_by"`
	ExpiresAt string  `json:"expires_at"`
	UsedAt    *string `json:"used_at,omitempty"`
	RevokedAt *string `json:"revoked_at,omitempty"`
	CreatedAt string  `json:"created_at"`
}

type inviteWithCreatorResponse struct {
	ID            string  `json:"id"`
	GroupID       string  `json:"group_id"`
	Code          string  `json:"code"`
	CreatedBy     string  `json:"created_by"`
	ExpiresAt     string  `json:"expires_at"`
	UsedAt        *string `json:"used_at,omitempty"`
	RevokedAt     *string `json:"revoked_at,omitempty"`
	CreatedAt     string  `json:"created_at"`
	CreatedByName string  `json:"created_by_name"`
}

type inviteWithGroupResponse struct {
	ID        string  `json:"id"`
	GroupID   string  `json:"group_id"`
	Code      string  `json:"code"`
	CreatedBy string  `json:"created_by"`
	ExpiresAt string  `json:"expires_at"`
	UsedAt    *string `json:"used_at,omitempty"`
	RevokedAt *string `json:"revoked_at,omitempty"`
	CreatedAt string  `json:"created_at"`
	GroupName string  `json:"group_name"`
}

func toInviteResponse(i invite.Invite) inviteResponse {
	r := inviteResponse{
		ID:        i.ID,
		GroupID:   i.GroupID,
		Code:      i.Code,
		CreatedBy: i.CreatedBy,
		ExpiresAt: i.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedAt: i.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if i.UsedAt != nil {
		s := i.UsedAt.Format("2006-01-02T15:04:05Z07:00")
		r.UsedAt = &s
	}
	if i.RevokedAt != nil {
		s := i.RevokedAt.Format("2006-01-02T15:04:05Z07:00")
		r.RevokedAt = &s
	}
	return r
}

func toInviteWithCreatorResponses(invites []invite.InviteWithCreator) []inviteWithCreatorResponse {
	res := make([]inviteWithCreatorResponse, len(invites))
	for i, inv := range invites {
		r := inviteWithCreatorResponse{
			ID:            inv.ID,
			GroupID:       inv.GroupID,
			Code:          inv.Code,
			CreatedBy:     inv.CreatedBy,
			ExpiresAt:     inv.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
			CreatedAt:     inv.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			CreatedByName: inv.CreatedByName,
		}
		if inv.UsedAt != nil {
			s := inv.UsedAt.Format("2006-01-02T15:04:05Z07:00")
			r.UsedAt = &s
		}
		if inv.RevokedAt != nil {
			s := inv.RevokedAt.Format("2006-01-02T15:04:05Z07:00")
			r.RevokedAt = &s
		}
		res[i] = r
	}
	return res
}

func toInviteWithGroupResponses(invites []invite.InviteWithGroup) []inviteWithGroupResponse {
	res := make([]inviteWithGroupResponse, len(invites))
	for i, inv := range invites {
		r := inviteWithGroupResponse{
			ID:        inv.ID,
			GroupID:   inv.GroupID,
			Code:      inv.Code,
			CreatedBy: inv.CreatedBy,
			ExpiresAt: inv.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
			CreatedAt: inv.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			GroupName: inv.GroupName,
		}
		if inv.UsedAt != nil {
			s := inv.UsedAt.Format("2006-01-02T15:04:05Z07:00")
			r.UsedAt = &s
		}
		if inv.RevokedAt != nil {
			s := inv.RevokedAt.Format("2006-01-02T15:04:05Z07:00")
			r.RevokedAt = &s
		}
		res[i] = r
	}
	return res
}
