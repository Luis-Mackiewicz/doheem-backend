package http

import (
	"encoding/json"
	"net/http"

	"doheem-backend/internal/domain"
)

type SplitTagHandler struct {
	svc *domain.SplitTagService
}

func NewSplitTagHandler(svc *domain.SplitTagService) *SplitTagHandler {
	return &SplitTagHandler{svc: svc}
}

func (h *SplitTagHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	groupID := r.PathValue("groupId")
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	tag, err := h.svc.Create(r.Context(), groupID, req.Name, userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toSplitTagResponse(tag))
}

func (h *SplitTagHandler) ListByGroup(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupId")
	tags, err := h.svc.ListByGroup(r.Context(), groupID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toSplitTagResponses(tags))
}

func (h *SplitTagHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	groupID := r.URL.Query().Get("group_id")
	if groupID == "" {
		respondError(w, http.StatusBadRequest, "group_id query param required")
		return
	}
	if err := h.svc.Delete(r.Context(), id, groupID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *SplitTagHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	members, err := h.svc.ListMembers(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toSplitTagMemberResponses(members))
}

func (h *SplitTagHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	member, err := h.svc.AddMember(r.Context(), id, req.UserID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toSplitTagMemberItemResponse(member))
}

func (h *SplitTagHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := r.PathValue("userId")
	if err := h.svc.RemoveMember(r.Context(), id, userID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
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
