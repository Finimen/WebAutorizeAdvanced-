package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSecretHandler struct {
	mock.Mock
}

func (m *MockSecretHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	w.Write([]byte("You are authorized! 🎉"))
}

func TestMiddelware_TableDriven(t *testing.T) {
	jwtKey := "test-secret-key"

	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		expectedCode   int
		expectedBody   string
		expectHandler  bool
		setupMockToken func() string
	}{
		{
			name: "Valid token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/secret", nil)
				token := generateValidToken(jwtKey, "testuser")
				req.Header.Set("Authorization", token)
				return req
			},
			expectedCode:  http.StatusOK,
			expectedBody:  "You are authorized! 🎉",
			expectHandler: true,
		},
		{
			name: "Missing token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/secret", nil)
				// Нет заголовка Authorization
				return req
			},
			expectedCode:  http.StatusUnauthorized,
			expectedBody:  "Missing token",
			expectHandler: false,
		},
		{
			name: "Invalid token format",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/secret", nil)
				req.Header.Set("Authorization", "invalid-token-format")
				return req
			},
			expectedCode:  http.StatusUnauthorized,
			expectedBody:  "Invalid token",
			expectHandler: false,
		},
		{
			name: "Expired token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/secret", nil)
				token := generateExpiredToken(jwtKey, "testuser")
				req.Header.Set("Authorization", token)
				return req
			},
			expectedCode:  http.StatusUnauthorized,
			expectedBody:  "Invalid token",
			expectHandler: false,
		},
		{
			name: "Wrong signing method",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/secret", nil)
				token := generateTokenWithWrongAlgorithm("testuser")
				req.Header.Set("Authorization", token)
				return req
			},
			expectedCode:  http.StatusUnauthorized,
			expectedBody:  "Invalid token",
			expectHandler: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHandler := &MockSecretHandler{}

			middleware := middelwareHandler(mockHandler.ServeHTTP, jwtKey)

			req := tt.setupRequest()
			rr := httptest.NewRecorder()

			if tt.expectHandler {
				mockHandler.On("ServeHTTP", rr, req).Return()
			}

			middleware.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)

			if tt.expectedBody != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedBody)
			}

			if tt.expectHandler {
				mockHandler.AssertCalled(t, "ServeHTTP", rr, req)
			} else {
				mockHandler.AssertNotCalled(t, "ServeHTTP", mock.Anything, mock.Anything)
			}
		})
	}
}

func generateValidToken(jwtKey, username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(jwtKey))
	return tokenString
}

func generateExpiredToken(jwtKey, username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(-time.Hour).Unix(), // Протухший токен
	})
	tokenString, _ := token.SignedString([]byte(jwtKey))
	return tokenString
}

func generateTokenWithWrongAlgorithm(username string) string {
	// Пытаемся использовать неправильный алгоритм (RS256 вместо HS256)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	if token == nil {
		return "invalid-token-with-wrong-alg"
	} else {
		return "tokenString"
	}
}
