package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	rdb      *redis.Client
	limit    int
	interval time.Duration
}

func NewRateLimiter(rdb *redis.Client, limit int, interval time.Duration) *RateLimiter {
	return &RateLimiter{rdb: rdb, limit: limit, interval: interval}
}

func (rl *RateLimiter) allow(ctx context.Context, key string) (bool, time.Duration, error) {
	now := time.Now().UnixNano()
	cutoff := now - rl.interval.Nanoseconds()
	redisKey := "ratelimit:" + key

	pipe := rl.rdb.Pipeline()
	pipe.ZRemRangeByScore(ctx, redisKey, "0", strconv.FormatInt(cutoff, 10))
	countCmd := pipe.ZCard(ctx, redisKey)
	pipe.ZAdd(ctx, redisKey, redis.Z{Score: float64(now), Member: now})
	pipe.Expire(ctx, redisKey, rl.interval)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, 0, fmt.Errorf("rate limiter: %w", err)
	}

	count, err := countCmd.Result()
	if err != nil {
		return false, 0, fmt.Errorf("rate limiter count: %w", err)
	}

	if int(count) > rl.limit {
		return false, rl.interval, nil
	}

	return true, 0, nil
}

func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return fwd
	}
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	return r.RemoteAddr
}

func RateLimitMiddleware(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := clientIP(r)
			if userID, ok := r.Context().Value(UserIDKey).(string); ok && userID != "" {
				key = "user:" + userID
			}

			allowed, retryAfter, err := rl.allow(r.Context(), key)
			if err != nil {
				slog.Error("rate limiter error", "error", err)
				next.ServeHTTP(w, r)
				return
			}
			if !allowed {
				sec := int(retryAfter.Seconds()) + 1
				w.Header().Set("Retry-After", strconv.Itoa(sec))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error":"too many requests"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
