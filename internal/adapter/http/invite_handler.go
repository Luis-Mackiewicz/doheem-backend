package http

import (
	"net/http"
	"time"

	"doheem-backend/internal/domain"
)

type InviteHandler struct {
	svc *domain.InviteService
}

func NewInviteHandler(svc *domain.InviteService) *InviteHandler {
	return &InviteHandler{svc: svc}
}

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
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toInviteResponse(invite))
}

func (h *InviteHandler) ListByGroup(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupId")
	invites, err := h.svc.ListByGroup(r.Context(), groupID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toInviteWithCreatorResponses(invites))
}

func (h *InviteHandler) ListPending(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	invites, err := h.svc.ListPendingByUser(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toInviteWithGroupResponses(invites))
}

func (h *InviteHandler) Accept(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	id := r.PathValue("id")
	if err := h.svc.Use(r.Context(), id, userID); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *InviteHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Revoke(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
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

func toInviteResponse(i domain.Invite) inviteResponse {
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

func toInviteWithCreatorResponses(invites []domain.InviteWithCreator) []inviteWithCreatorResponse {
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

func toInviteWithGroupResponses(invites []domain.InviteWithGroup) []inviteWithGroupResponse {
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
