package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type IPasswordHasher interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword([]byte, []byte) error
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type BcryptHasher struct {
}

func (b *BcryptHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func (b *BcryptHasher) CompareHashAndPassword(storedPaswsord []byte, userPassword []byte) error {
	err := bcrypt.CompareHashAndPassword(storedPaswsord, userPassword)

	if err != nil {
		return err
	}

	return nil
}

type RegisterHandler struct {
	UserRepo IRepository
	Hasher   IPasswordHasher
}

func (h *RegisterHandler) registerHandler(w http.ResponseWriter, r *http.Request) {
	cxt := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Use post", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	hashedPassword, err := h.Hasher.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Hashing password error", http.StatusInternalServerError)
		return
	}

	err = h.UserRepo.CreateUser(cxt, user.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Username already exist", http.StatusBadRequest)
		return
	}

	w.Write([]byte("User registrated successfuly"))
}
