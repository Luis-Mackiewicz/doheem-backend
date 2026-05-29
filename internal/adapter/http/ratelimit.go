package http

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	windows  map[string]*slidingWindow
	limit    int
	interval time.Duration
}

type slidingWindow struct {
	timestamps []time.Time
}

func NewRateLimiter(limit int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		windows:  make(map[string]*slidingWindow),
		limit:    limit,
		interval: interval,
	}
	go rl.cleanup(5 * interval)
	return rl
}

func (rl *RateLimiter) allow(key string) (bool, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.interval)

	w, ok := rl.windows[key]
	if !ok {
		w = &slidingWindow{}
		rl.windows[key] = w
	}

	filtered := w.timestamps[:0]
	for _, t := range w.timestamps {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	w.timestamps = filtered

	if len(w.timestamps) >= rl.limit {
		oldest := w.timestamps[0]
		retryAfter := rl.interval - now.Sub(oldest)
		return false, retryAfter
	}

	w.timestamps = append(w.timestamps, now)
	return true, 0
}

func (rl *RateLimiter) cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		cutoff := time.Now().Add(-rl.interval)
		for key, w := range rl.windows {
			filtered := w.timestamps[:0]
			for _, t := range w.timestamps {
				if t.After(cutoff) {
					filtered = append(filtered, t)
				}
			}
			if len(filtered) == 0 {
				delete(rl.windows, key)
			} else {
				w.timestamps = filtered
			}
		}
		rl.mu.Unlock()
	}
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

			allowed, retryAfter := rl.allow(key)
			if !allowed {
				sec := int(retryAfter.Seconds()) + 1
				w.Header().Set("Retry-After", itoa(sec))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error":"too many requests"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}
