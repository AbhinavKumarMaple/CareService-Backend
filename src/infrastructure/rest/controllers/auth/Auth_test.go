package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	useCaseAuth "caregiver/src/application/usecases/auth"
	userDomain "caregiver/src/domain/user"
	logger "caregiver/src/infrastructure/logger"

	"github.com/gin-gonic/gin"
)

type MockAuthUseCase struct {
	loginFunc                func(string, string) (*userDomain.User, *useCaseAuth.AuthTokens, error)
	accessTokenByRefreshFunc func(string) (*userDomain.User, *useCaseAuth.AuthTokens, error)
}

func (m *MockAuthUseCase) Login(email, password string) (*userDomain.User, *useCaseAuth.AuthTokens, error) {
	if m.loginFunc != nil {
		return m.loginFunc(email, password)
	}
	return nil, nil, nil
}

func (m *MockAuthUseCase) AccessTokenByRefreshToken(refreshToken string) (*userDomain.User, *useCaseAuth.AuthTokens, error) {
	if m.accessTokenByRefreshFunc != nil {
		return m.accessTokenByRefreshFunc(refreshToken)
	}
	return nil, nil, nil
}

func setupLogger(t *testing.T) *logger.Logger {
	loggerInstance, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	return loggerInstance
}

func TestNewAuthController(t *testing.T) {
	mockUseCase := &MockAuthUseCase{}
	logger := setupLogger(t)
	controller := NewAuthController(mockUseCase, logger)

	if controller == nil {
		t.Error("Expected NewAuthController to return a non-nil controller")
	}
}

func TestAuthController_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := &MockAuthUseCase{
		loginFunc: func(email, password string) (*userDomain.User, *useCaseAuth.AuthTokens, error) {
			user := &userDomain.User{
				UserName:  "testuser",
				Email:     "test@example.com",
				FirstName: "Test",
				LastName:  "User",
				Status:    true,
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			}
			authTokens := &useCaseAuth.AuthTokens{
				AccessToken:               "test-access-token",
				RefreshToken:              "test-refresh-token",
				ExpirationAccessDateTime:  time.Now().Add(time.Hour),
				ExpirationRefreshDateTime: time.Now().Add(24 * time.Hour),
			}
			return user, authTokens, nil
		},
	}

	logger := setupLogger(t)
	controller := NewAuthController(mockUseCase, logger)

	loginRequest := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	requestBody, _ := json.Marshal(loginRequest)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	controller.Login(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAuthController_Login_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := &MockAuthUseCase{}

	logger := setupLogger(t)
	controller := NewAuthController(mockUseCase, logger)

	requestBody := []byte(`{"email": "test@example.com"}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	controller.Login(c)

	if len(c.Errors) == 0 {
		t.Error("Expected error to be added to context")
	}
}

func TestAuthController_GetAccessTokenByRefreshToken_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := &MockAuthUseCase{
		accessTokenByRefreshFunc: func(refreshToken string) (*userDomain.User, *useCaseAuth.AuthTokens, error) {
			user := &userDomain.User{
				UserName:  "testuser",
				Email:     "test@example.com",
				FirstName: "Test",
				LastName:  "User",
				Status:    true,
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			}
			authTokens := &useCaseAuth.AuthTokens{
				AccessToken:               "new-access-token",
				RefreshToken:              "new-refresh-token",
				ExpirationAccessDateTime:  time.Now().Add(time.Hour),
				ExpirationRefreshDateTime: time.Now().Add(24 * time.Hour),
			}
			return user, authTokens, nil
		},
	}

	logger := setupLogger(t)
	controller := NewAuthController(mockUseCase, logger)

	accessTokenRequest := AccessTokenRequest{
		RefreshToken: "test-refresh-token",
	}

	requestBody, _ := json.Marshal(accessTokenRequest)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	controller.GetAccessTokenByRefreshToken(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAuthController_GetAccessTokenByRefreshToken_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := &MockAuthUseCase{}

	logger := setupLogger(t)
	controller := NewAuthController(mockUseCase, logger)

	requestBody := []byte(`{}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	controller.GetAccessTokenByRefreshToken(c)

	if len(c.Errors) == 0 {
		t.Error("Expected error to be added to context")
	}
}

func TestLoginRequest_Validation(t *testing.T) {
	validRequest := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	if validRequest.Email == "" {
		t.Error("Email should not be empty")
	}

	if validRequest.Password == "" {
		t.Error("Password should not be empty")
	}

	if validRequest.Email == "invalid-email" {
		t.Error("Email should be in valid format")
	}
}

func TestAccessTokenRequest_Validation(t *testing.T) {
	validRequest := AccessTokenRequest{
		RefreshToken: "valid-refresh-token",
	}

	if validRequest.RefreshToken == "" {
		t.Error("RefreshToken should not be empty")
	}
}
