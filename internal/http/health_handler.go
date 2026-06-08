package http

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	pool *pgxpool.Pool
	rdb  *redis.Client
}

func NewHealthHandler(pool *pgxpool.Pool, rdb *redis.Client) *HealthHandler {
	return &HealthHandler{pool: pool, rdb: rdb}
}

type healthCheck struct {
	Status  string            `json:"status"`
	Checks  map[string]string `json:"checks"`
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	result := healthCheck{
		Checks: make(map[string]string),
	}
	statusCode := http.StatusOK

	// Database
	if err := h.pool.Ping(ctx); err != nil {
		result.Checks["database"] = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	} else {
		result.Checks["database"] = "healthy"
	}

	// Redis
	if err := h.rdb.Ping(ctx).Err(); err != nil {
		result.Checks["redis"] = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	} else {
		result.Checks["redis"] = "healthy"
	}

	if statusCode == http.StatusOK {
		result.Status = "ok"
	} else {
		result.Status = "degraded"
	}

	respondJSON(w, statusCode, result)
}
