package main

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterHandler_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  map[string]interface{}
		setupMocks   func(*MockUserRepository, *MockPasswordHasher)
		expectedCode int
		expectedBody string
	}{
		{
			name: "Succesful regestration",
			requestBody: map[string]interface{}{
				"username": "validuser",
				"password": "password123",
			},
			setupMocks: func(mur *MockUserRepository, mph *MockPasswordHasher) {
				mur.On("CreateUser", "validuser", "hashed_password").Return(nil)
				mph.On("GenerateFromPassword", []byte("password123"), bcrypt.DefaultCost).Return([]byte("hashed_password"), nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: "User registrated successfuly",
		},
		{
			name: "empty username",
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
			name: "empty password",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"password": "",
			},
			setupMocks:   func(mockRepo *MockUserRepository, mockHasher *MockPasswordHasher) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid input",
		},
		{
			name: "username already exists",
			requestBody: map[string]interface{}{
				"username": "existinguser",
				"password": "password123",
			},
			setupMocks: func(mockRepo *MockUserRepository, mockHasher *MockPasswordHasher) {
				mockHasher.On("GenerateFromPassword", mock.Anything, mock.Anything).
					Return([]byte("hashed_password"), nil)
				mockRepo.On("CreateUser", "existinguser", "hashed_password").
					Return(errors.New("ErrUserExists"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Username already exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository := new(MockUserRepository)
			mockHasher := new(MockPasswordHasher)

			tt.setupMocks(mockRepository, mockHasher)

			handler := RegisterHandler{
				UserRepo: mockRepository,
				Hasher:   mockHasher,
			}

			req := createTestRequest(http.MethodPost, "/register", tt.requestBody)

			rr := executeHandler(handler.registerHandler, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)

			mockRepository.AssertExpectations(t)
			mockHasher.AssertExpectations(t)
		})
	}
}
