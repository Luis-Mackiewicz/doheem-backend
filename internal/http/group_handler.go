package http

import (
	"net/http"

	"doheem-backend/internal/group"
)

type GroupHandler struct {
	svc *group.GroupService
}

func NewGroupHandler(svc *group.GroupService) *GroupHandler {
	return &GroupHandler{svc: svc}
}

func (h *GroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	var req struct {
		Name string `json:"name" validate:"required"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	g, err := h.svc.Create(r.Context(), group.CreateGroupParams{
		Name: req.Name,
	}, userID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusCreated, toGroupResponse(g))
}

func (h *GroupHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	groups, err := h.svc.ListByUser(r.Context(), userID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	limit, offset := parsePagination(r)
	items, total := paginate(toGroupResponses(groups), limit, offset)
	respondJSON(w, http.StatusOK, paginatedResponse{Data: items, Total: total})
}

func (h *GroupHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("id")
	members, err := h.svc.ListMembers(r.Context(), groupID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	limit, offset := parsePagination(r)
	items, total := paginate(toGroupMemberResponses(members), limit, offset)
	respondJSON(w, http.StatusOK, paginatedResponse{Data: items, Total: total})
}

func (h *GroupHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	group, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, toGroupResponse(group))
}

func (h *GroupHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		Name        *string  `json:"name,omitempty"`
		Description *string  `json:"description,omitempty"`
		MonthlyFee  *float64 `json:"monthly_fee,omitempty"`
		Cnpj        *string  `json:"cnpj,omitempty"`
		Cep         *string  `json:"cep,omitempty"`
		PhotoURL    *string  `json:"photo_url,omitempty"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	g, err := h.svc.Update(r.Context(), id, group.UpdateGroupParams{
		Name:        req.Name,
		Description: req.Description,
		MonthlyFee:  req.MonthlyFee,
		Cnpj:        req.Cnpj,
		Cep:         req.Cep,
		PhotoURL:    req.PhotoURL,
	})
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, toGroupResponse(g))
}

func (h *GroupHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("id")
	var req struct {
		UserID  string `json:"user_id" validate:"required"`
		IsAdmin bool   `json:"is_admin"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	member, err := h.svc.AddMember(r.Context(), groupID, req.UserID, req.IsAdmin)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusCreated, toGroupMemberResponse(member))
}

func (h *GroupHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("id")
	userID := r.PathValue("userId")
	var req struct {
		IsAdmin bool `json:"is_admin"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	member, err := h.svc.UpdateMemberRole(r.Context(), groupID, userID, req.IsAdmin)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, toGroupMemberResponse(member))
}

func (h *GroupHandler) Join(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	groupID := r.PathValue("id")
	if err := h.svc.Join(r.Context(), groupID, userID); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *GroupHandler) RegenerateInvite(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("id")
	token, err := h.svc.RegenerateInviteToken(r.Context(), groupID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]*string{"invite_token": token})
}

func (h *GroupHandler) GetInviteToken(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("id")
	token, err := h.svc.GetInviteToken(r.Context(), groupID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]*string{"invite_token": token})
}

func (h *GroupHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("id")
	userID := r.PathValue("userId")
	if err := h.svc.RemoveMember(r.Context(), groupID, userID); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type groupResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	MonthlyFee  float64 `json:"monthly_fee"`
	Cnpj        string  `json:"cnpj"`
	Cep         string  `json:"cep"`
	PhotoURL    *string `json:"photo_url,omitempty"`
	InviteToken *string `json:"invite_token,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type groupMemberResponse struct {
	ID       string `json:"id"`
	GroupID  string `json:"group_id"`
	UserID   string `json:"user_id"`
	IsAdmin  bool   `json:"is_admin"`
	JoinedAt string `json:"joined_at"`
}

type groupMemberWithUserResponse struct {
	ID        string  `json:"id"`
	GroupID   string  `json:"group_id"`
	UserID    string  `json:"user_id"`
	IsAdmin   bool    `json:"is_admin"`
	JoinedAt  string  `json:"joined_at"`
	UserName  string  `json:"user_name"`
	UserEmail string  `json:"user_email"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

func toGroupResponse(g group.Group) groupResponse {
	return groupResponse{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		MonthlyFee:  g.MonthlyFee,
		Cnpj:        g.Cnpj,
		Cep:         g.Cep,
		PhotoURL:    g.PhotoURL,
		InviteToken: g.InviteToken,
		CreatedAt:   g.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   g.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toGroupResponses(groups []group.Group) []groupResponse {
	res := make([]groupResponse, len(groups))
	for i, g := range groups {
		res[i] = toGroupResponse(g)
	}
	return res
}

func toGroupMemberResponse(m group.GroupMember) groupMemberResponse {
	return groupMemberResponse{
		ID:       m.ID,
		GroupID:  m.GroupID,
		UserID:   m.UserID,
		IsAdmin:  m.IsAdmin,
		JoinedAt: m.JoinedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toGroupMemberResponses(members []group.GroupMemberWithUser) []groupMemberWithUserResponse {
	res := make([]groupMemberWithUserResponse, len(members))
	for i, m := range members {
		res[i] = groupMemberWithUserResponse{
			ID:        m.ID,
			GroupID:   m.GroupID,
			UserID:    m.UserID,
			IsAdmin:   m.IsAdmin,
			JoinedAt:  m.JoinedAt.Format("2006-01-02T15:04:05Z07:00"),
			UserName:  m.UserName,
			UserEmail: m.UserEmail,
			AvatarURL: m.AvatarURL,
		}
	}
	return res
}
