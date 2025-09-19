package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func secretHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You are authorized! 🎉"))
}

func middelwareHandler(next http.HandlerFunc, jwtKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		r = r.WithContext(ctx)

		tokenStr := r.Header.Get("Authorization")

		log.Printf("Access attempt from %s to %s", getClientIP(r), r.URL.Path)

		if tokenStr == "" {
			log.Printf("Missing token from %s", getClientIP(r))
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(jwtKey), nil
		})

		if err != nil || !token.Valid {
			log.Printf("Invalid token from %s: %v", getClientIP(r), err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		log.Printf("User logged in")
		next.ServeHTTP(w, r)
	}
}
