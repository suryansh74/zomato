package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	authModel "github.com/suryansh74/zomato/services/auth-service/internal/models"
	"github.com/suryansh74/zomato/services/auth-service/internal/services"
	"github.com/suryansh74/zomato/services/shared/middleware"
	sharedModels "github.com/suryansh74/zomato/services/shared/models"
	"github.com/suryansh74/zomato/services/shared/token"
	"golang.org/x/oauth2"
)

type mockAuthRepository struct {
	mock.Mock
}

func (m *mockAuthRepository) FindByEmail(ctx context.Context, email string) (*authModel.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.User), args.Error(1)
}

func (m *mockAuthRepository) FindByID(ctx context.Context, id string) (*authModel.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.User), args.Error(1)
}

func (m *mockAuthRepository) Create(ctx context.Context, user *authModel.User) (*authModel.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.User), args.Error(1)
}

func (m *mockAuthRepository) UpdateRole(ctx context.Context, role authModel.Role, id string) (*authModel.User, error) {
	args := m.Called(ctx, role, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.User), args.Error(1)
}

type mockTokenMaker struct {
	mock.Mock
}

func (m *mockTokenMaker) CreateToken(user *sharedModels.TokenUser, duration time.Duration) (string, error) {
	args := m.Called(user, duration)
	return args.String(0), args.Error(1)
}

func (m *mockTokenMaker) VerifyToken(tokenStr string) (*token.Payload, error) {
	args := m.Called(tokenStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*token.Payload), args.Error(1)
}

func createTestHandler(mockRepo *mockAuthRepository, mockMaker *mockTokenMaker) *AuthHandler {
	svc := services.NewAuthService(mockRepo)

	oauthConfig := &oauth2.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/api/auth/google/callback",
		Scopes:       []string{"email", "profile"},
	}

	return NewAuthHandler(
		svc,
		mockMaker,
		24*time.Hour,
		oauthConfig,
		true,
		"http://localhost:5173",
	)
}

func TestAuthHandler_CheckHealth(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	mockMaker := new(mockTokenMaker)
	handler := createTestHandler(mockRepo, mockMaker)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/health", nil)
	w := httptest.NewRecorder()

	handler.CheckHealth(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "auth-service is healthy", data["message"])
}

func TestAuthHandler_Login(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	mockMaker := new(mockTokenMaker)
	handler := createTestHandler(mockRepo, mockMaker)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/login", nil)
	w := httptest.NewRecorder()

	handler.Login(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)

	location := resp.Header.Get("Location")
	assert.Contains(t, location, "access_type=offline")
	assert.Contains(t, location, "state=login")
}

func TestAuthHandler_GoogleCallback_MissingCode(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	mockMaker := new(mockTokenMaker)
	handler := createTestHandler(mockRepo, mockMaker)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/google/callback", nil)
	w := httptest.NewRecorder()

	handler.GoogleCallback(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "missing code", data["error"])
}

func TestAuthHandler_Profile_Success(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	mockMaker := new(mockTokenMaker)
	handler := createTestHandler(mockRepo, mockMaker)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/profile", nil)

	payload := &token.Payload{
		User: &sharedModels.TokenUser{
			ID:    "test-user-id",
			Name:  "Test User",
			Email: "test@example.com",
			Image: "https://example.com/img.png",
			Role:  "customer",
		},
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	ctx := context.WithValue(req.Context(), middleware.UserContextKey, payload)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.Profile(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "profile successful", data["message"])

	payloadData, ok := data["payload"].(map[string]interface{})
	require.True(t, ok)

	userData, ok := payloadData["user"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "Test User", userData["name"])
	assert.Equal(t, "test@example.com", userData["email"])
}

func TestAuthHandler_AddRole_Success(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	mockMaker := new(mockTokenMaker)

	updatedUser := &authModel.User{
		Name:  "Test User",
		Email: "test@example.com",
		Image: "https://example.com/img.png",
		Role:  "restaurant_owner",
	}

	mockRepo.On("UpdateRole", mock.Anything, mock.AnythingOfType("models.Role"), mock.AnythingOfType("string")).Return(updatedUser, nil)
	mockMaker.On("CreateToken", mock.AnythingOfType("*models.TokenUser"), mock.AnythingOfType("time.Duration")).Return("new-mock-token", nil)

	handler := createTestHandler(mockRepo, mockMaker)

	body := `{"role": "restaurant_owner"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/add_role", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	payload := &token.Payload{
		User: &sharedModels.TokenUser{
			ID:    "test-user-id",
			Name:  "Test User",
			Email: "test@example.com",
			Image: "https://example.com/img.png",
			Role:  "customer",
		},
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	ctx := context.WithValue(req.Context(), middleware.UserContextKey, payload)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.AddRole(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "role updated successfully", data["message"])

	cookies := w.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "session_token" {
			sessionCookie = c
			break
		}
	}
	assert.NotNil(t, sessionCookie)
	assert.Equal(t, "new-mock-token", sessionCookie.Value)

	mockRepo.AssertExpectations(t)
	mockMaker.AssertExpectations(t)
}

func TestAuthHandler_AddRole_InvalidBody(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	mockMaker := new(mockTokenMaker)
	handler := createTestHandler(mockRepo, mockMaker)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/add_role", strings.NewReader("invalid-json"))
	req.Header.Set("Content-Type", "application/json")

	payload := &token.Payload{
		User: &sharedModels.TokenUser{
			ID:    "test-user-id",
			Name:  "Test User",
			Email: "test@example.com",
			Image: "https://example.com/img.png",
			Role:  "customer",
		},
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	ctx := context.WithValue(req.Context(), middleware.UserContextKey, payload)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.AddRole(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "invalid request body", data["error"])
}

func TestAuthHandler_AddRole_InvalidRole(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	mockMaker := new(mockTokenMaker)
	handler := createTestHandler(mockRepo, mockMaker)

	body := `{"role": "invalid_role"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/add_role", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	payload := &token.Payload{
		User: &sharedModels.TokenUser{
			ID:    "test-user-id",
			Name:  "Test User",
			Email: "test@example.com",
			Image: "https://example.com/img.png",
			Role:  "customer",
		},
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	ctx := context.WithValue(req.Context(), middleware.UserContextKey, payload)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.AddRole(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestAuthHandler_Logout(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	mockMaker := new(mockTokenMaker)
	handler := createTestHandler(mockRepo, mockMaker)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	w := httptest.NewRecorder()

	handler.Logout(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "logged out", data["message"])

	cookies := w.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "session_token" {
			sessionCookie = c
			break
		}
	}

	assert.NotNil(t, sessionCookie)
	assert.Equal(t, "", sessionCookie.Value)
	assert.Equal(t, -1, sessionCookie.MaxAge)
	assert.True(t, sessionCookie.HttpOnly)
}

func TestAuthHandler_setSessionCookie(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	mockMaker := new(mockTokenMaker)
	handler := createTestHandler(mockRepo, mockMaker)

	w := httptest.NewRecorder()

	handler.setSessionCookie(w, "test-token-value")

	cookies := w.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "session_token" {
			sessionCookie = c
			break
		}
	}

	require.NotNil(t, sessionCookie)
	assert.Equal(t, "test-token-value", sessionCookie.Value)
	assert.Equal(t, "/", sessionCookie.Path)
	assert.True(t, sessionCookie.HttpOnly)
	assert.Equal(t, http.SameSiteLaxMode, sessionCookie.SameSite)
	assert.False(t, sessionCookie.Secure)
}

func TestAuthHandler_setSessionCookie_Production(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	mockMaker := new(mockTokenMaker)

	oauthConfig := &oauth2.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/api/auth/google/callback",
		Scopes:       []string{"email", "profile"},
	}

	svc := services.NewAuthService(mockRepo)
	handler := NewAuthHandler(
		svc,
		mockMaker,
		24*time.Hour,
		oauthConfig,
		false,
		"http://localhost:5173",
	)

	w := httptest.NewRecorder()
	handler.setSessionCookie(w, "test-token-value")

	cookies := w.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "session_token" {
			sessionCookie = c
			break
		}
	}

	require.NotNil(t, sessionCookie)
	assert.True(t, sessionCookie.Secure)
}

func TestAuthHandler_AddRole_EmptyRole(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	mockMaker := new(mockTokenMaker)
	handler := createTestHandler(mockRepo, mockMaker)

	body := `{"role": ""}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/add_role", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	payload := &token.Payload{
		User: &sharedModels.TokenUser{
			ID:    "test-user-id",
			Name:  "Test User",
			Email: "test@example.com",
			Image: "https://example.com/img.png",
			Role:  "customer",
		},
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	ctx := context.WithValue(req.Context(), middleware.UserContextKey, payload)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.AddRole(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestAuthHandler_AddRole_ValidRoles(t *testing.T) {
	validRoles := []string{"customer", "restaurant_owner", "rider"}

	for _, roleName := range validRoles {
		t.Run(roleName, func(t *testing.T) {
			mockRepo := new(mockAuthRepository)
			mockMaker := new(mockTokenMaker)

			updatedUser := &authModel.User{
				Name:  "Test User",
				Email: "test@example.com",
				Image: "https://example.com/img.png",
				Role:  roleName,
			}

			mockRepo.On("UpdateRole", mock.Anything, mock.AnythingOfType("models.Role"), mock.AnythingOfType("string")).Return(updatedUser, nil)
			mockMaker.On("CreateToken", mock.AnythingOfType("*models.TokenUser"), mock.AnythingOfType("time.Duration")).Return("new-mock-token", nil)

			handler := createTestHandler(mockRepo, mockMaker)

			body := `{"role": "` + roleName + `"}`
			req := httptest.NewRequest(http.MethodPost, "/api/auth/add_role", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			payload := &token.Payload{
				User: &sharedModels.TokenUser{
					ID:    "test-user-id",
					Name:  "Test User",
					Email: "test@example.com",
					Image: "https://example.com/img.png",
					Role:  "customer",
				},
				IssuedAt:  time.Now(),
				ExpiredAt: time.Now().Add(24 * time.Hour),
			}

			ctx := context.WithValue(req.Context(), middleware.UserContextKey, payload)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			handler.AddRole(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			mockRepo.AssertExpectations(t)
			mockMaker.AssertExpectations(t)
		})
	}
}
