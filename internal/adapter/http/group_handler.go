package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"doheem-backend/internal/domain"
)

type GroupHandler struct {
	svc *domain.GroupService
}

func NewGroupHandler(svc *domain.GroupService) *GroupHandler {
	return &GroupHandler{svc: svc}
}

func (h *GroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	var req struct {
		Name     string `json:"name"`
		Currency string `json:"currency"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	group, err := h.svc.Create(r.Context(), domain.CreateGroupParams{
		Name:     req.Name,
		Currency: req.Currency,
	}, userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toGroupResponse(group))
}

func (h *GroupHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	groups, err := h.svc.ListByUser(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toGroupResponses(groups))
}

func (h *GroupHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	group, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toGroupResponse(group))
}

func (h *GroupHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		Name     *string `json:"name,omitempty"`
		Currency *string `json:"currency,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	group, err := h.svc.Update(r.Context(), id, domain.UpdateGroupParams{
		Name:     req.Name,
		Currency: req.Currency,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toGroupResponse(group))
}

func (h *GroupHandler) SoftDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.SoftDelete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *GroupHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Deactivate(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *GroupHandler) Activate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Activate(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *GroupHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("id")
	members, err := h.svc.ListMembers(r.Context(), groupID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toGroupMemberResponses(members))
}

func (h *GroupHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("id")
	var req struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	member, err := h.svc.AddMember(r.Context(), groupID, req.UserID, req.Role)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toGroupMemberResponse(member))
}

func (h *GroupHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("id")
	userID := r.PathValue("userId")
	var req struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	member, err := h.svc.UpdateMemberRole(r.Context(), groupID, userID, req.Role)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toGroupMemberResponse(member))
}

func (h *GroupHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("id")
	userID := r.PathValue("userId")
	if err := h.svc.RemoveMember(r.Context(), groupID, userID); err != nil {
		if errors.Is(err, domain.ErrCannotRemoveOwner) {
			respondError(w, http.StatusForbidden, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type groupResponse struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Currency      string  `json:"currency"`
	IsActive      bool    `json:"is_active"`
	InactiveSince *string `json:"inactive_since,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type groupMemberResponse struct {
	ID       string `json:"id"`
	GroupID  string `json:"group_id"`
	UserID   string `json:"user_id"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
	JoinedAt string `json:"joined_at"`
}

type groupMemberWithUserResponse struct {
	ID        string  `json:"id"`
	GroupID   string  `json:"group_id"`
	UserID    string  `json:"user_id"`
	Role      string  `json:"role"`
	IsActive  bool    `json:"is_active"`
	JoinedAt  string  `json:"joined_at"`
	UserName  string  `json:"user_name"`
	UserEmail string  `json:"user_email"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

func toGroupResponse(g domain.Group) groupResponse {
	r := groupResponse{
		ID:        g.ID,
		Name:      g.Name,
		Currency:  g.Currency,
		IsActive:  g.IsActive,
		CreatedAt: g.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: g.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if g.InactiveSince != nil {
		s := g.InactiveSince.Format("2006-01-02T15:04:05Z07:00")
		r.InactiveSince = &s
	}
	return r
}

func toGroupResponses(groups []domain.Group) []groupResponse {
	res := make([]groupResponse, len(groups))
	for i, g := range groups {
		res[i] = toGroupResponse(g)
	}
	return res
}

func toGroupMemberResponse(m domain.GroupMember) groupMemberResponse {
	return groupMemberResponse{
		ID:       m.ID,
		GroupID:  m.GroupID,
		UserID:   m.UserID,
		Role:     m.Role,
		IsActive: m.IsActive,
		JoinedAt: m.JoinedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toGroupMemberResponses(members []domain.GroupMemberWithUser) []groupMemberWithUserResponse {
	res := make([]groupMemberWithUserResponse, len(members))
	for i, m := range members {
		res[i] = groupMemberWithUserResponse{
			ID:        m.ID,
			GroupID:   m.GroupID,
			UserID:    m.UserID,
			Role:      m.Role,
			IsActive:  m.IsActive,
			JoinedAt:  m.JoinedAt.Format("2006-01-02T15:04:05Z07:00"),
			UserName:  m.UserName,
			UserEmail: m.UserEmail,
			AvatarURL: m.AvatarURL,
		}
	}
	return res
}
