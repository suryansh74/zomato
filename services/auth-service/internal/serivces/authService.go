package services

import (
	"context"

	"github.com/suryansh74/zomato/services/auth-service/apperr"
	"github.com/suryansh74/zomato/services/auth-service/internal/models"
	"github.com/suryansh74/zomato/services/auth-service/internal/repositories"
)

type AuthService struct {
	repo repositories.AuthRepository
}

func NewAuthService(repo repositories.AuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

// LoginOrCreate if user already exists return user else create user
// ================================================================
func (a *AuthService) LoginOrCreate(ctx context.Context, LoginRequest *models.LoginRequest) (*models.LoginRequest, error) {
	// check if user already exists
	_, err := a.repo.FindByEmail(ctx, LoginRequest.Email)

	// create user if not exists
	if err == apperr.ErrUserNotFound {
		_, err = a.repo.Create(ctx, LoginRequest)
		if err != nil {
			return nil, apperr.ErrInternalServer
		}
		return LoginRequest, nil
	}

	if err != nil {
		return nil, apperr.ErrInternalServer
	}
	return LoginRequest, nil
}
