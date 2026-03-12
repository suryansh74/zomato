package serivces

import (
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

// LoginOrCreate
func (a *AuthService) LoginOrCreate(LoginRequest *models.LoginRequest) (*models.LoginRequest, error) {
	// check if user already exists
	_, err := a.repo.FindByEmail(LoginRequest.Email)
	if err != nil {
		// create user if not exists
		if err == apperr.ErrUserNotFound {
			_, err = a.repo.Create(LoginRequest)
			if err != nil {
				return nil, apperr.ErrInternalServer
			}
			return LoginRequest, nil
		}
		return nil, apperr.ErrInternalServer
	}
	return LoginRequest, nil
}
