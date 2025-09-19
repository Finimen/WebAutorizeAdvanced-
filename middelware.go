package main

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func secretHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You are authorized! 🎉"))
}

func middelwareHandler(next http.HandlerFunc, jwtKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")

		if tokenStr == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
