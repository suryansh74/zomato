package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/suryansh74/zomato/services/auth-service/apperr"
	"github.com/suryansh74/zomato/services/auth-service/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthRepository) UpdateRole(ctx context.Context, role models.Role, id string) (*models.User, error) {
	args := m.Called(ctx, role, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestMockAuthRepository_FindByEmail_Success(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	expectedUser := &models.User{
		ID:        bson.NewObjectID(),
		Name:      "Test User",
		Email:     "test@example.com",
		Image:     "https://example.com/image.png",
		Role:      "customer",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(expectedUser, nil)

	user, err := mockRepo.FindByEmail(context.Background(), "test@example.com")

	require.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestMockAuthRepository_FindByEmail_NotFound(t *testing.T) {
	mockRepo := new(MockAuthRepository)

	mockRepo.On("FindByEmail", mock.Anything, "unknown@example.com").Return((*models.User)(nil), apperr.ErrUserNotFound)

	user, err := mockRepo.FindByEmail(context.Background(), "unknown@example.com")

	assert.Nil(t, user)
	assert.Equal(t, apperr.ErrUserNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestMockAuthRepository_FindByEmail_Error(t *testing.T) {
	mockRepo := new(MockAuthRepository)

	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").Return((*models.User)(nil), apperr.ErrInternalServer)

	user, err := mockRepo.FindByEmail(context.Background(), "test@example.com")

	assert.Nil(t, user)
	assert.Equal(t, apperr.ErrInternalServer, err)
	mockRepo.AssertExpectations(t)
}

func TestMockAuthRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	userID := bson.NewObjectID()
	expectedUser := &models.User{
		ID:        userID,
		Name:      "Test User",
		Email:     "test@example.com",
		Image:     "https://example.com/image.png",
		Role:      "admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByID", mock.Anything, userID.Hex()).Return(expectedUser, nil)

	user, err := mockRepo.FindByID(context.Background(), userID.Hex())

	require.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestMockAuthRepository_FindByID_NotFound(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	userID := bson.NewObjectID()

	mockRepo.On("FindByID", mock.Anything, userID.Hex()).Return((*models.User)(nil), apperr.ErrUserNotFound)

	user, err := mockRepo.FindByID(context.Background(), userID.Hex())

	assert.Nil(t, user)
	assert.Equal(t, apperr.ErrUserNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestMockAuthRepository_Create_Success(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	newUser := &models.User{
		Name:  "New User",
		Email: "new@example.com",
		Image: "https://example.com/image.png",
	}
	createdUser := &models.User{
		ID:        bson.NewObjectID(),
		Name:      newUser.Name,
		Email:     newUser.Email,
		Image:     newUser.Image,
		Role:      "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("Create", mock.Anything, newUser).Return(createdUser, nil)

	user, err := mockRepo.Create(context.Background(), newUser)

	require.NoError(t, err)
	assert.Equal(t, createdUser, user)
	mockRepo.AssertExpectations(t)
}

func TestMockAuthRepository_Create_Error(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	newUser := &models.User{
		Name:  "New User",
		Email: "new@example.com",
		Image: "https://example.com/image.png",
	}

	mockRepo.On("Create", mock.Anything, newUser).Return((*models.User)(nil), apperr.ErrInternalServer)

	user, err := mockRepo.Create(context.Background(), newUser)

	assert.Nil(t, user)
	assert.Equal(t, apperr.ErrInternalServer, err)
	mockRepo.AssertExpectations(t)
}

func TestMockAuthRepository_UpdateRole_Success(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	userID := bson.NewObjectID()
	role := models.Role{Role: "admin"}
	updatedUser := &models.User{
		ID:        userID,
		Name:      "Test User",
		Email:     "test@example.com",
		Image:     "https://example.com/image.png",
		Role:      "admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("UpdateRole", mock.Anything, role, userID.Hex()).Return(updatedUser, nil)

	user, err := mockRepo.UpdateRole(context.Background(), role, userID.Hex())

	require.NoError(t, err)
	assert.Equal(t, "admin", user.Role)
	assert.Equal(t, updatedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestMockAuthRepository_UpdateRole_Error(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	userID := bson.NewObjectID()
	role := models.Role{Role: "admin"}

	mockRepo.On("UpdateRole", mock.Anything, role, userID.Hex()).Return((*models.User)(nil), apperr.ErrInternalServer)

	user, err := mockRepo.UpdateRole(context.Background(), role, userID.Hex())

	assert.Nil(t, user)
	assert.Equal(t, apperr.ErrInternalServer, err)
	mockRepo.AssertExpectations(t)
}

func TestMockAuthRepository_FindByEmail_ContextCancellation(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").Return((*models.User)(nil), context.Canceled)

	user, err := mockRepo.FindByEmail(ctx, "test@example.com")

	assert.Nil(t, user)
	assert.True(t, errors.Is(err, context.Canceled))
	mockRepo.AssertExpectations(t)
}

func TestMockAuthRepository_Create_EmptyEmail(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	newUser := &models.User{
		Name:  "New User",
		Email: "",
		Image: "https://example.com/image.png",
	}

	mockRepo.On("Create", mock.Anything, newUser).Return((*models.User)(nil), apperr.ErrInternalServer)

	user, err := mockRepo.Create(context.Background(), newUser)

	assert.Nil(t, user)
	assert.Equal(t, apperr.ErrInternalServer, err)
	mockRepo.AssertExpectations(t)
}

func TestMockAuthRepository_UpdateRole_InvalidRole(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	userID := bson.NewObjectID()
	role := models.Role{Role: "invalid_role"}

	mockRepo.On("UpdateRole", mock.Anything, role, userID.Hex()).Return((*models.User)(nil), apperr.ErrInternalServer)

	user, err := mockRepo.UpdateRole(context.Background(), role, userID.Hex())

	assert.Nil(t, user)
	assert.Equal(t, apperr.ErrInternalServer, err)
	mockRepo.AssertExpectations(t)
}

func TestMockAuthRepository_FindByID_InvalidHex(t *testing.T) {
	mockRepo := new(MockAuthRepository)

	mockRepo.On("FindByID", mock.Anything, "invalid-hex").Return((*models.User)(nil), apperr.ErrInternalServer)

	user, err := mockRepo.FindByID(context.Background(), "invalid-hex")

	assert.Nil(t, user)
	assert.Equal(t, apperr.ErrInternalServer, err)
	mockRepo.AssertExpectations(t)
}

func TestMockAuthRepository_UpdateRole_RoleValues(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	userID := bson.NewObjectID()
	roles := []string{"customer", "restaurant_owner", "rider"}

	for _, roleName := range roles {
		t.Run(roleName, func(t *testing.T) {
			role := models.Role{Role: roleName}
			updatedUser := &models.User{
				ID:    userID,
				Name:  "Test User",
				Email: "test@example.com",
				Image: "https://example.com/image.png",
				Role:  roleName,
			}

			mockRepo.On("UpdateRole", mock.Anything, role, userID.Hex()).Return(updatedUser, nil)

			user, err := mockRepo.UpdateRole(context.Background(), role, userID.Hex())

			require.NoError(t, err)
			assert.Equal(t, roleName, user.Role)
			mockRepo.AssertExpectations(t)
		})
	}
}
