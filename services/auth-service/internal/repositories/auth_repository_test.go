package repositories_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/suryansh74/zomato/services/auth-service/apperr"
	"github.com/suryansh74/zomato/services/auth-service/internal/models"
	"github.com/suryansh74/zomato/services/auth-service/internal/repositories"
)

// ─── Suite Setup ─────────────────────────────────────────────────────────────

type AuthRepositoryTestSuite struct {
	suite.Suite
	container testcontainers.Container
	client    *mongo.Client
	repo      repositories.AuthRepository
	ctx       context.Context
}

func (s *AuthRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mongo:7",
		ExposedPorts: []string{"27017/tcp"},
		Tmpfs:        map[string]string{"/data/db": "rw"}, // ✅ fixes permissions issue
		WaitingFor:   wait.ForListeningPort("27017/tcp"),
	}

	container, err := testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(s.T(), err)
	s.container = container

	host, err := container.Host(s.ctx)
	require.NoError(s.T(), err)

	port, err := container.MappedPort(s.ctx, "27017")
	require.NoError(s.T(), err)

	uri := fmt.Sprintf("mongodb://%s:%s", host, port.Port())
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	require.NoError(s.T(), err)
	s.client = client

	s.repo = repositories.NewAuthRepository(client, "testdb", "users")
}

func (s *AuthRepositoryTestSuite) TearDownSuite() {
	s.client.Disconnect(s.ctx)
	s.container.Terminate(s.ctx)
}

// clears collection before each test so tests don't affect each other
func (s *AuthRepositoryTestSuite) SetupTest() {
	s.client.Database("testdb").Collection("users").Drop(s.ctx)
}

func TestAuthRepository(t *testing.T) {
	suite.Run(t, new(AuthRepositoryTestSuite))
}

// ─── FindByEmail Tests ────────────────────────────────────────────────────────

func (s *AuthRepositoryTestSuite) Test_FindByEmail_WhenUserExists_ReturnsUser() {
	// arrange
	req := &models.LoginRequest{
		Name:  "Suryansh",
		Email: "suryansh@gmail.com",
		Image: "https://image.com/suryansh.jpg",
	}
	_, err := s.repo.Create(s.ctx, req)
	require.NoError(s.T(), err)

	// act
	result, err := s.repo.FindByEmail(s.ctx, "suryansh@gmail.com")

	// assert
	require.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "suryansh@gmail.com", result["email"])
	assert.Equal(s.T(), "Suryansh", result["name"])
}

func (s *AuthRepositoryTestSuite) Test_FindByEmail_WhenUserDoesNotExist_ReturnsErrUserNotFound() {
	// act
	result, err := s.repo.FindByEmail(s.ctx, "ghost@gmail.com")

	// assert
	assert.Nil(s.T(), result)
	assert.ErrorIs(s.T(), err, apperr.ErrUserNotFound)
}

func (s *AuthRepositoryTestSuite) Test_FindByEmail_WhenEmailCaseIsDifferent_ReturnsNotFound() {
	// arrange
	req := &models.LoginRequest{
		Name:  "Suryansh",
		Email: "suryansh@gmail.com",
		Image: "https://image.com/suryansh.jpg",
	}
	_, err := s.repo.Create(s.ctx, req)
	require.NoError(s.T(), err)

	// act — search with uppercase, MongoDB is case sensitive by default
	result, err := s.repo.FindByEmail(s.ctx, "SURYANSH@GMAIL.COM")

	// assert
	assert.Nil(s.T(), result)
	assert.ErrorIs(s.T(), err, apperr.ErrUserNotFound)
}

// ─── Create Tests ─────────────────────────────────────────────────────────────

func (s *AuthRepositoryTestSuite) Test_Create_WhenValidRequest_InsertsUser() {
	// arrange
	req := &models.LoginRequest{
		Name:  "Suryansh",
		Email: "suryansh@gmail.com",
		Image: "https://image.com/suryansh.jpg",
	}

	// act
	result, err := s.repo.Create(s.ctx, req)

	// assert
	require.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.NotNil(s.T(), result.InsertedID)
}

func (s *AuthRepositoryTestSuite) Test_Create_WhenUserInserted_CanBeFoundByEmail() {
	// arrange
	req := &models.LoginRequest{
		Name:  "Suryansh",
		Email: "suryansh@gmail.com",
		Image: "https://image.com/suryansh.jpg",
	}

	// act
	_, err := s.repo.Create(s.ctx, req)
	require.NoError(s.T(), err)

	// assert — verify it actually persisted in DB
	found, err := s.repo.FindByEmail(s.ctx, req.Email)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), req.Email, found["email"])
	assert.Equal(s.T(), req.Name, found["name"])
}

func (s *AuthRepositoryTestSuite) Test_Create_WhenUserInserted_HasTimestamps() {
	before := time.Now().Add(-time.Second)
	req := &models.LoginRequest{
		Name:  "Suryansh",
		Email: "suryansh@gmail.com",
		Image: "https://image.com/suryansh.jpg",
	}

	_, err := s.repo.Create(s.ctx, req)
	require.NoError(s.T(), err)

	found, err := s.repo.FindByEmail(s.ctx, req.Email)
	require.NoError(s.T(), err)

	// ✅ bson.M returns timestamps as primitive.DateTime, convert to time.Time
	createdAtRaw, ok := found["created_at"]
	require.True(s.T(), ok, "created_at field should exist")

	var createdAt time.Time
	switch v := createdAtRaw.(type) {
	case bson.DateTime: // ✅ v2 type
		createdAt = v.Time()
	case time.Time:
		createdAt = v
	default:
		s.T().Fatalf("unexpected type for created_at: %T", createdAtRaw)
	}

	assert.True(s.T(), createdAt.After(before), "created_at should be after test start")
}
