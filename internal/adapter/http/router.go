package http

import (
	"net/http"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	health := &HealthHandler{}
	mux.HandleFunc("GET /", health.HealthCheck)
	return mux
}
