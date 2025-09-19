package main

import (
	"net/http"
	"sort"
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
	windowStart := now.Add(-rl.window)

	// Получаем таймстемпы для IP
	timestamps, exists := rl.requests[ip]
	if !exists {
		timestamps = []time.Time{}
	}

	// Находим индекс первого таймстемпа, который ещё в окне
	// (все таймстемпы после этого индекса актуальны)
	firstValidIndex := sort.Search(len(timestamps), func(i int) bool {
		return timestamps[i].After(windowStart)
	})
	validTimestamps := timestamps[firstValidIndex:]

	// Проверяем лимит
	if len(validTimestamps) >= rl.maxRequests {
		return false
	}

	// Разрешаем запрос, добавляем текущую метку и сохраняем только актуальные
	validTimestamps = append(validTimestamps, now)
	rl.requests[ip] = validTimestamps // Сохраняем уже очищенный срез

	return true
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
