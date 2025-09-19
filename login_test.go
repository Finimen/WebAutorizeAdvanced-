package main

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginHandler_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  map[string]interface{}
		setupMocks   func(*MockUserRepository, *MockPasswordHasher)
		expectedCode int
		expectedBody string
		checkToken   bool
	}{
		{
			name: "Successful login",
			requestBody: map[string]interface{}{
				"username": "validuser",
				"password": "correctpassword",
			},
			setupMocks: func(mur *MockUserRepository, mph *MockPasswordHasher) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
				mur.On("GetUserByUsername", "validuser").Return(string(hashedPassword), nil)

				mph.On("CompareHashAndPassword", []byte(hashedPassword), []byte("correctpassword")).Return(nil)
			},
			expectedCode: http.StatusOK,
			checkToken:   true,
		},
		{
			name: "Wrong password",
			requestBody: map[string]interface{}{
				"username": "validuser",
				"password": "wrongpassword",
			},
			setupMocks: func(mur *MockUserRepository, mph *MockPasswordHasher) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
				mur.On("GetUserByUsername", "validuser").Return(string(hashedPassword), nil)

				mph.On("CompareHashAndPassword", []byte(hashedPassword), []byte("wrongpassword")).
					Return(bcrypt.ErrMismatchedHashAndPassword)
			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Invalid creditans",
		},
		{
			name: "User not found",
			requestBody: map[string]interface{}{
				"username": "nonexistent",
				"password": "anypassword",
			},
			setupMocks: func(mur *MockUserRepository, mph *MockPasswordHasher) {
				mur.On("GetUserByUsername", "nonexistent").Return("", sql.ErrNoRows)
			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Invalid credentials",
		},
		{
			name: "Empty username",
			requestBody: map[string]interface{}{
				"username": "",
				"password": "password123",
			},
			setupMocks: func(mur *MockUserRepository, mph *MockPasswordHasher) {
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid input",
		},
		{
			name: "Empty password",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"password": "",
			},
			setupMocks: func(mur *MockUserRepository, mph *MockPasswordHasher) {
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid input",
		},
		{
			name: "Invalid JSON",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"password": 12345,
			},
			setupMocks: func(mur *MockUserRepository, mph *MockPasswordHasher) {
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository := &MockUserRepository{}
			mockHasher := &MockPasswordHasher{}

			tt.setupMocks(mockRepository, mockHasher)

			handler := LoginHandler{
				Repo:   mockRepository,
				Hasher: mockHasher,
				JwtKey: []byte("test-secret-key"),
			}

			req := createTestRequest(http.MethodPost, "/login", tt.requestBody)
			rr := executeHandler(handler.loginHandler, req)

			assert.Equal(t, tt.expectedCode, rr.Code)

			if tt.expectedBody != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedBody)
			}

			if tt.checkToken {
				tokenString := rr.Body.String()
				assert.NotEmpty(t, tokenString)

				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					return []byte("test-secret-key"), nil
				})

				assert.NoError(t, err)
				assert.True(t, token.Valid)

				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					assert.Equal(t, "validuser", claims["username"])
					assert.NotEmpty(t, claims["exp"])
				}
			}

			mockRepository.AssertExpectations(t)
			mockHasher.AssertExpectations(t)
		})
	}
}
