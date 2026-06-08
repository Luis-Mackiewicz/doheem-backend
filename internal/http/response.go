package http

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type paginatedResponse struct {
	Data  any `json:"data"`
	Total int `json:"total"`
}

func parsePagination(r *http.Request) (limit, offset int) {
	limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return
}

func paginate[T any](items []T, limit, offset int) ([]T, int) {
	total := len(items)
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return items[offset:end], total
}

func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		json.NewEncoder(w).Encode(v)
	}
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}

func respondValidationError(w http.ResponseWriter, errs []validationError) {
	respondJSON(w, http.StatusBadRequest, map[string]any{
		"error":  "validation failed",
		"fields": errs,
	})
}
