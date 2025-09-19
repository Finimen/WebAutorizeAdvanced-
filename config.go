package main

import "time"

type Config struct {
	JWTKey     string
	Port       string
	DBPath     string
	RateLimit  int           // Максимальное количество запросов
	RateWindow time.Duration // Временное окно для rate limiting
}

// Дефолтная конфигурация
func DefaultConfig() Config {
	return Config{
		JWTKey:     "secretKey",
		Port:       "8888",
		DBPath:     "./users.db",
		RateLimit:  100,         // 100 запросов
		RateWindow: time.Minute, // в течение 1 минуты
	}
}
