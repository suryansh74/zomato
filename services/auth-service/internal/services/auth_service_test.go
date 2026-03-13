package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/suryansh74/zomato/services/auth-service/apperr"
	"github.com/suryansh74/zomato/services/auth-service/internal/models"
	"github.com/suryansh74/zomato/services/auth-service/internal/services"
)

// ─── Mock Repository ─────────────────────────────────────────────────────────
// we mock repo here because service tests should NOT hit real DB
// service tests only care about business logic

type mockAuthRepo struct {
	mock.Mock
}

func (m *mockAuthRepo) FindByEmail(ctx context.Context, email string) (bson.M, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(bson.M), args.Error(1)
}

func (m *mockAuthRepo) Create(ctx context.Context, req *models.LoginRequest) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

// ─── Suite Setup ─────────────────────────────────────────────────────────────

type AuthServiceTestSuite struct {
	suite.Suite
	mockRepo *mockAuthRepo
	svc      *services.AuthService
	ctx      context.Context
}

func (s *AuthServiceTestSuite) SetupTest() {
	s.mockRepo = new(mockAuthRepo)
	s.svc = services.NewAuthService(s.mockRepo)
	s.ctx = context.Background()
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

// ─── LoginOrCreate Tests ──────────────────────────────────────────────────────

func (s *AuthServiceTestSuite) Test_LoginOrCreate_WhenUserExists_ReturnsUserWithoutCreating() {
	// arrange
	req := &models.LoginRequest{
		Name:  "Suryansh",
		Email: "suryansh@gmail.com",
		Image: "https://image.com/suryansh.jpg",
	}
	existingUser := bson.M{"email": req.Email, "name": req.Name}

	s.mockRepo.On("FindByEmail", s.ctx, req.Email).Return(existingUser, nil)

	// act
	result, err := s.svc.LoginOrCreate(s.ctx, req)

	// assert
	require.NoError(s.T(), err)
	assert.Equal(s.T(), req, result)
	// Create should NEVER be called when user exists
	s.mockRepo.AssertNotCalled(s.T(), "Create")
}

func (s *AuthServiceTestSuite) Test_LoginOrCreate_WhenUserNotFound_CreatesAndReturnsUser() {
	// arrange
	req := &models.LoginRequest{
		Name:  "Suryansh",
		Email: "suryansh@gmail.com",
		Image: "https://image.com/suryansh.jpg",
	}

	s.mockRepo.On("FindByEmail", s.ctx, req.Email).Return(nil, apperr.ErrUserNotFound)
	s.mockRepo.On("Create", s.ctx, req).Return(&mongo.InsertOneResult{}, nil)

	// act
	result, err := s.svc.LoginOrCreate(s.ctx, req)

	// assert
	require.NoError(s.T(), err)
	assert.Equal(s.T(), req, result)
	s.mockRepo.AssertCalled(s.T(), "Create", s.ctx, req)
}

func (s *AuthServiceTestSuite) Test_LoginOrCreate_WhenUserNotFound_AndCreateFails_ReturnsInternalError() {
	// arrange
	req := &models.LoginRequest{
		Name:  "Suryansh",
		Email: "suryansh@gmail.com",
		Image: "https://image.com/suryansh.jpg",
	}

	s.mockRepo.On("FindByEmail", s.ctx, req.Email).Return(nil, apperr.ErrUserNotFound)
	s.mockRepo.On("Create", s.ctx, req).Return(nil, apperr.ErrInternalServer)

	// act
	result, err := s.svc.LoginOrCreate(s.ctx, req)

	// assert
	assert.Nil(s.T(), result)
	assert.ErrorIs(s.T(), err, apperr.ErrInternalServer)
}

func (s *AuthServiceTestSuite) Test_LoginOrCreate_WhenFindByEmailFails_ReturnsInternalError() {
	// arrange
	req := &models.LoginRequest{
		Name:  "Suryansh",
		Email: "suryansh@gmail.com",
		Image: "https://image.com/suryansh.jpg",
	}

	s.mockRepo.On("FindByEmail", s.ctx, req.Email).Return(nil, apperr.ErrInternalServer)

	// act
	result, err := s.svc.LoginOrCreate(s.ctx, req)

	// assert
	assert.Nil(s.T(), result)
	assert.ErrorIs(s.T(), err, apperr.ErrInternalServer)
	// Create should NEVER be called when FindByEmail has unexpected error
	s.mockRepo.AssertNotCalled(s.T(), "Create")
}
