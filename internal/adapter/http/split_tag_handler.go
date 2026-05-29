package http

import (
	"net/http"

	"doheem-backend/internal/domain"
)

type SplitTagHandler struct {
	svc *domain.SplitTagService
}

func NewSplitTagHandler(svc *domain.SplitTagService) *SplitTagHandler {
	return &SplitTagHandler{svc: svc}
}

// Create creates a new split tag in a group
// @Summary Create a split tag
// @Description Create a new split tag for a group
// @Tags Split Tags
// @Accept json
// @Produce json
// @Param groupId path string true "Group ID"
// @Param request body object{name=string} true "Tag name"
// @Success 201 {object} splitTagResponse
// @Failure 400 {object} map[string]any "Validation error"
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/groups/{groupId}/split-tags [post]
// @Security BearerAuth
func (h *SplitTagHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	groupID := r.PathValue("groupId")
	var req struct {
		Name string `json:"name" validate:"required"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	tag, err := h.svc.Create(r.Context(), groupID, req.Name, userID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusCreated, toSplitTagResponse(tag))
}

// ListByGroup lists split tags in a group
// @Summary List split tags by group
// @Description List all split tags for a specific group
// @Tags Split Tags
// @Accept json
// @Produce json
// @Param groupId path string true "Group ID"
// @Success 200 {array} splitTagResponse
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Group not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/groups/{groupId}/split-tags [get]
// @Security BearerAuth
func (h *SplitTagHandler) ListByGroup(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupId")
	tags, err := h.svc.ListByGroup(r.Context(), groupID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, toSplitTagResponses(tags))
}

// Delete deletes a split tag
// @Summary Delete a split tag
// @Description Permanently delete a split tag
// @Tags Split Tags
// @Accept json
// @Produce json
// @Param id path string true "Split Tag ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]any "group_id query param required"
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Split tag not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/split-tags/{id} [delete]
// @Security BearerAuth
func (h *SplitTagHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	groupID := r.URL.Query().Get("group_id")
	if groupID == "" {
		respondError(w, http.StatusBadRequest, "group_id query param required")
		return
	}
	if err := h.svc.Delete(r.Context(), id, groupID); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListMembers lists members of a split tag
// @Summary List split tag members
// @Description List all members assigned to a split tag
// @Tags Split Tags
// @Accept json
// @Produce json
// @Param id path string true "Split Tag ID"
// @Success 200 {array} splitTagMemberResponse
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Split tag not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/split-tags/{id}/members [get]
// @Security BearerAuth
func (h *SplitTagHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	members, err := h.svc.ListMembers(r.Context(), id)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, toSplitTagMemberResponses(members))
}

// AddMember adds a member to a split tag
// @Summary Add split tag member
// @Description Add a member to a split tag
// @Tags Split Tags
// @Accept json
// @Produce json
// @Param id path string true "Split Tag ID"
// @Param request body object{user_id=string} true "User ID"
// @Success 201 {object} splitTagMemberItemResponse
// @Failure 400 {object} map[string]any "Validation error"
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Split tag not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/split-tags/{id}/members [post]
// @Security BearerAuth
func (h *SplitTagHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		UserID string `json:"user_id" validate:"required"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	member, err := h.svc.AddMember(r.Context(), id, req.UserID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusCreated, toSplitTagMemberItemResponse(member))
}

// RemoveMember removes a member from a split tag
// @Summary Remove split tag member
// @Description Remove a member from a split tag
// @Tags Split Tags
// @Accept json
// @Produce json
// @Param id path string true "Split Tag ID"
// @Param userId path string true "User ID"
// @Success 204 {object} nil
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/split-tags/{id}/members/{userId} [delete]
// @Security BearerAuth
func (h *SplitTagHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := r.PathValue("userId")
	if err := h.svc.RemoveMember(r.Context(), id, userID); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type splitTagResponse struct {
	ID        string `json:"id"`
	GroupID   string `json:"group_id"`
	Name      string `json:"name"`
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`
}

type splitTagMemberResponse struct {
	ID         string `json:"id"`
	SplitTagID string `json:"split_tag_id"`
	UserID     string `json:"user_id"`
	UserName   string `json:"user_name"`
}

type splitTagMemberItemResponse struct {
	ID         string `json:"id"`
	SplitTagID string `json:"split_tag_id"`
	UserID     string `json:"user_id"`
}

func toSplitTagResponse(t domain.SplitTag) splitTagResponse {
	return splitTagResponse{
		ID:        t.ID,
		GroupID:   t.GroupID,
		Name:      t.Name,
		CreatedBy: t.CreatedBy,
		CreatedAt: t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toSplitTagResponses(tags []domain.SplitTag) []splitTagResponse {
	res := make([]splitTagResponse, len(tags))
	for i, t := range tags {
		res[i] = toSplitTagResponse(t)
	}
	return res
}

func toSplitTagMemberResponses(members []domain.SplitTagMemberWithUser) []splitTagMemberResponse {
	res := make([]splitTagMemberResponse, len(members))
	for i, m := range members {
		res[i] = splitTagMemberResponse{
			ID:         m.ID,
			SplitTagID: m.SplitTagID,
			UserID:     m.UserID,
			UserName:   m.UserName,
		}
	}
	return res
}

func toSplitTagMemberItemResponse(m domain.SplitTagMember) splitTagMemberItemResponse {
	return splitTagMemberItemResponse{
		ID:         m.ID,
		SplitTagID: m.SplitTagID,
		UserID:     m.UserID,
	}
}
