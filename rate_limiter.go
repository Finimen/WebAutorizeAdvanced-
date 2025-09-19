package main

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu          sync.Mutex
	requests    map[string][]time.Time
	maxRequests int
	window      time.Duration
}

func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests:    make(map[string][]time.Time),
		maxRequests: maxRequests,
		window:      window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	rl.requests[ip] = rl.cleanup(rl.requests[ip], now)

	if len(rl.requests[ip]) >= rl.maxRequests {
		return false
	}

	rl.requests[ip] = append(rl.requests[ip], now)
	return true
}

func (rl *RateLimiter) cleanup(request []time.Time, now time.Time) []time.Time {
	var cleaned []time.Time
	for _, t := range request {
		if now.Sub(t) <= rl.window {
			cleaned = append(cleaned, t)
		}
	}

	return cleaned
}

func RateLimitMiddleware(limiter *RateLimiter, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		if !limiter.Allow(ip) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}
