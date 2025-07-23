package utils

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
	}

	return limiter
}

func (rl *RateLimiter) cleanupOldEntries() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for ip, limiter := range rl.limiters {
		if limiter.Tokens() == float64(rl.burst) {
			// Limiter hasn't been used recently, remove it
			delete(rl.limiters, ip)
		}
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	// Start cleanup routine
	go func() {
		ticker := time.NewTicker(time.Minute * 5)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanupOldEntries()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getRealIP(r)
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			RespondWithError(w, http.StatusTooManyRequests, APIError{
				Detail: "rate limit exceeded",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getRealIP(r *http.Request) string {
	// Try various headers to get the real IP
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if firstIP := net.ParseIP(ip); firstIP != nil {
			return firstIP.String()
		}
	}

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		if realIP := net.ParseIP(ip); realIP != nil {
			return realIP.String()
		}
	}

	// Fall back to remote address
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
