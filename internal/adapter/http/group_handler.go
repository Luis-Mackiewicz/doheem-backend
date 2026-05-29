package http

import (
	"net/http"

	"doheem-backend/internal/domain"
)

type GroupHandler struct {
	svc *domain.GroupService
}

func NewGroupHandler(svc *domain.GroupService) *GroupHandler {
	return &GroupHandler{svc: svc}
}

// Create creates a new group
// @Summary Create a group
// @Description Create a new group with the authenticated user as owner
// @Tags Groups
// @Accept json
// @Produce json
// @Param request body object{name=string,currency=string} true "Group details"
// @Success 201 {object} groupResponse
// @Failure 400 {object} map[string]any "Validation error"
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/groups [post]
// @Security BearerAuth
func (h *GroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	var req struct {
		Name     string `json:"name"     validate:"required"`
		Currency string `json:"currency" validate:"required,len=3"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	group, err := h.svc.Create(r.Context(), domain.CreateGroupParams{
		Name:     req.Name,
		Currency: req.Currency,
	}, userID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusCreated, toGroupResponse(group))
}

// List lists all groups for the authenticated user
// @Summary List groups
// @Description List all groups the authenticated user is a member of
// @Tags Groups
// @Accept json
// @Produce json
// @Success 200 {array} groupResponse
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/groups [get]
// @Security BearerAuth
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

// GetByID gets a group by ID
// @Summary Get group by ID
// @Description Get a group by its unique identifier
// @Tags Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} groupResponse
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Group not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/groups/{id} [get]
// @Security BearerAuth
func (h *GroupHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	group, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, toGroupResponse(group))
}

// Update updates a group
// @Summary Update a group
// @Description Update an existing group's details
// @Tags Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body object{name=string,currency=string} true "Group update details"
// @Success 200 {object} groupResponse
// @Failure 400 {object} map[string]any "Validation error"
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Group not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/groups/{id} [put]
// @Security BearerAuth
func (h *GroupHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		Name     *string `json:"name,omitempty"`
		Currency *string `json:"currency,omitempty" validate:"omitempty,len=3"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	group, err := h.svc.Update(r.Context(), id, domain.UpdateGroupParams{
		Name:     req.Name,
		Currency: req.Currency,
	})
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, toGroupResponse(group))
}

// SoftDelete soft-deletes a group
// @Summary Soft delete a group
// @Description Soft delete (mark as deleted) a group
// @Tags Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Success 204 {object} nil
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Group not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/groups/{id} [delete]
// @Security BearerAuth
func (h *GroupHandler) SoftDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.SoftDelete(r.Context(), id); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Deactivate deactivates a group
// @Summary Deactivate a group
// @Description Deactivate a group, making it temporarily inactive
// @Tags Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Success 204 {object} nil
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Group not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/groups/{id}/deactivate [patch]
// @Security BearerAuth
func (h *GroupHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Deactivate(r.Context(), id); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Activate activates a group
// @Summary Activate a group
// @Description Reactivate a previously deactivated group
// @Tags Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Success 204 {object} nil
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Group not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/groups/{id}/activate [patch]
// @Security BearerAuth
func (h *GroupHandler) Activate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Activate(r.Context(), id); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListMembers lists members of a group
// @Summary List group members
// @Description List all members of a specific group
// @Tags Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {array} groupMemberWithUserResponse
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Group not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/groups/{id}/members [get]
// @Security BearerAuth
// AddMember adds a member to a group
// @Summary Add group member
// @Description Add a new member to a group with a specific role
// @Tags Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body object{user_id=string,role=string} true "Member details"
// @Success 201 {object} groupMemberResponse
// @Failure 400 {object} map[string]any "Validation error"
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Group not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/groups/{id}/members [post]
// @Security BearerAuth
func (h *GroupHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("id")
	var req struct {
		UserID string `json:"user_id" validate:"required"`
		Role   string `json:"role"    validate:"required,oneof=owner admin member"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	member, err := h.svc.AddMember(r.Context(), groupID, req.UserID, req.Role)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusCreated, toGroupMemberResponse(member))
}

// UpdateMemberRole updates a member's role in a group
// @Summary Update member role
// @Description Update the role of a member in a group
// @Tags Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param userId path string true "User ID"
// @Param request body object{role=string} true "New role"
// @Success 200 {object} groupMemberResponse
// @Failure 400 {object} map[string]any "Validation error"
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/groups/{id}/members/{userId} [put]
// @Security BearerAuth
func (h *GroupHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("id")
	userID := r.PathValue("userId")
	var req struct {
		Role string `json:"role" validate:"required,oneof=owner admin member"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	member, err := h.svc.UpdateMemberRole(r.Context(), groupID, userID, req.Role)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, toGroupMemberResponse(member))
}

// RemoveMember removes a member from a group
// @Summary Remove group member
// @Description Remove a member from a group
// @Tags Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param userId path string true "User ID"
// @Success 204 {object} nil
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/groups/{id}/members/{userId} [delete]
// @Security BearerAuth
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
