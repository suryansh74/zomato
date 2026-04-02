package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/suryansh74/zomato/services/auth-service/apperr"
	"github.com/suryansh74/zomato/services/auth-service/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type mockAuthRepository struct {
	mock.Mock
}

func (m *mockAuthRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockAuthRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockAuthRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockAuthRepository) UpdateRole(ctx context.Context, role models.Role, id string) (*models.User, error) {
	args := m.Called(ctx, role, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestAuthService_LoginOrCreate_ExistingUser(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	svc := NewAuthService(mockRepo)

	existingUser := &models.User{
		ID:        bson.NewObjectID(),
		Name:      "Existing User",
		Email:     "existing@example.com",
		Image:     "https://example.com/img.png",
		Role:      "customer",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
	}

	mockRepo.On("FindByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)

	req := &models.LoginRequest{
		Name:  "Existing User",
		Email: "existing@example.com",
		Image: "https://example.com/img.png",
	}

	user, err := svc.LoginOrCreate(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, existingUser, user)
	assert.Equal(t, "customer", user.Role)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_LoginOrCreate_NewUser(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	svc := NewAuthService(mockRepo)

	newUser := &models.User{
		ID:        bson.NewObjectID(),
		Name:      "New User",
		Email:     "new@example.com",
		Image:     "https://example.com/img.png",
		Role:      "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByEmail", mock.Anything, "new@example.com").Return((*models.User)(nil), apperr.ErrUserNotFound)
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
		return u.Name == "New User" && u.Email == "new@example.com" && u.Image == "https://example.com/img.png"
	})).Return(newUser, nil)

	req := &models.LoginRequest{
		Name:  "New User",
		Email: "new@example.com",
		Image: "https://example.com/img.png",
	}

	user, err := svc.LoginOrCreate(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, newUser, user)
	assert.Empty(t, user.Role)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_LoginOrCreate_FindByEmailError(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	svc := NewAuthService(mockRepo)

	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").Return((*models.User)(nil), apperr.ErrInternalServer)

	req := &models.LoginRequest{
		Name:  "Test User",
		Email: "test@example.com",
		Image: "https://example.com/img.png",
	}

	user, err := svc.LoginOrCreate(context.Background(), req)

	assert.Nil(t, user)
	assert.Equal(t, apperr.ErrInternalServer, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_LoginOrCreate_CreateError(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	svc := NewAuthService(mockRepo)

	mockRepo.On("FindByEmail", mock.Anything, "new@example.com").Return((*models.User)(nil), apperr.ErrUserNotFound)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return((*models.User)(nil), apperr.ErrInternalServer)

	req := &models.LoginRequest{
		Name:  "New User",
		Email: "new@example.com",
		Image: "https://example.com/img.png",
	}

	user, err := svc.LoginOrCreate(context.Background(), req)

	assert.Nil(t, user)
	assert.Equal(t, apperr.ErrInternalServer, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_UpdateRole_Success(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	svc := NewAuthService(mockRepo)

	userID := bson.NewObjectID()
	userIDHex := userID.Hex()
	role := models.Role{Role: "admin"}
	updatedUser := &models.User{
		ID:        userID,
		Name:      "Test User",
		Email:     "test@example.com",
		Image:     "https://example.com/img.png",
		Role:      "admin",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("UpdateRole", mock.Anything, role, userIDHex).Return(updatedUser, nil)

	result, err := svc.UpdateRole(context.Background(), role, userIDHex)

	require.NoError(t, err)
	assert.Equal(t, updatedUser, result)
	assert.Equal(t, "admin", result.Role)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_UpdateRole_Error(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	svc := NewAuthService(mockRepo)

	userID := bson.NewObjectID()
	userIDHex := userID.Hex()
	role := models.Role{Role: "customer"}

	mockRepo.On("UpdateRole", mock.Anything, role, userIDHex).Return((*models.User)(nil), apperr.ErrInternalServer)

	result, err := svc.UpdateRole(context.Background(), role, userIDHex)

	assert.Nil(t, result)
	assert.Equal(t, apperr.ErrInternalServer, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_UpdateRole_AllRoles(t *testing.T) {
	roles := []string{"customer", "restaurant_owner", "rider", "admin"}

	for _, roleName := range roles {
		t.Run(roleName, func(t *testing.T) {
			mockRepo := new(mockAuthRepository)
			svc := NewAuthService(mockRepo)

			objID := bson.NewObjectID()
			userID := objID.Hex()
			role := models.Role{Role: roleName}
			updatedUser := &models.User{
				ID:    objID,
				Name:  "Test User",
				Email: "test@example.com",
				Image: "https://example.com/img.png",
				Role:  roleName,
			}

			mockRepo.On("UpdateRole", mock.Anything, role, userID).Return(updatedUser, nil)

			result, err := svc.UpdateRole(context.Background(), role, userID)

			require.NoError(t, err)
			assert.Equal(t, roleName, result.Role)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_LoginOrCreate_ContextCancellation(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	svc := NewAuthService(mockRepo)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").Return((*models.User)(nil), context.Canceled)

	req := &models.LoginRequest{
		Name:  "Test User",
		Email: "test@example.com",
		Image: "https://example.com/img.png",
	}

	user, err := svc.LoginOrCreate(ctx, req)

	assert.Nil(t, user)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_NewAuthService(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	svc := NewAuthService(mockRepo)

	assert.NotNil(t, svc)
	assert.Equal(t, mockRepo, svc.repo)
}

func TestAuthService_UpdateRole_InvalidUserID(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	svc := NewAuthService(mockRepo)

	invalidID := "not-a-valid-object-id"
	role := models.Role{Role: "admin"}

	mockRepo.On("UpdateRole", mock.Anything, role, invalidID).Return((*models.User)(nil), apperr.ErrInternalServer)

	result, err := svc.UpdateRole(context.Background(), role, invalidID)

	assert.Nil(t, result)
	assert.Equal(t, apperr.ErrInternalServer, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_LoginOrCreate_EmptyEmail(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	svc := NewAuthService(mockRepo)

	mockRepo.On("FindByEmail", mock.Anything, "").Return((*models.User)(nil), apperr.ErrUserNotFound)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return((*models.User)(nil), apperr.ErrInternalServer)

	req := &models.LoginRequest{
		Name:  "Test User",
		Email: "",
		Image: "https://example.com/img.png",
	}

	user, err := svc.LoginOrCreate(context.Background(), req)

	assert.Nil(t, user)
	assert.Equal(t, apperr.ErrInternalServer, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_UpdateRole_ContextCancellation(t *testing.T) {
	mockRepo := new(mockAuthRepository)
	svc := NewAuthService(mockRepo)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	userID := bson.NewObjectID().Hex()
	role := models.Role{Role: "admin"}

	mockRepo.On("UpdateRole", mock.Anything, role, userID).Return((*models.User)(nil), context.Canceled)

	result, err := svc.UpdateRole(ctx, role, userID)

	assert.Nil(t, result)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
