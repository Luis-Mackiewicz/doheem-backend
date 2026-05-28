package http

import (
	"fmt"
	"net/http"
)

type HealthHandler struct{}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Doheem server is running! 🚀")
}
