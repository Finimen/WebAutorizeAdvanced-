package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginHandler struct {
	Repo   SQLRepository
	Hasher BcryptHasher
	JwtKey []byte
}

func (l *LoginHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Use post", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	storedPassword, err := l.Repo.GetUserByUsername(user.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = l.Hasher.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid creditans", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	tokenstring, err := token.SignedString(l.JwtKey)

	if err != nil {
		http.Error(w, "Token generating error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenstring))
}
