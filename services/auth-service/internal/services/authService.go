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
	return &AuthService{repo: repo}
}

func (a *AuthService) LoginOrCreate(ctx context.Context, req *models.LoginRequest) (*models.User, error) {
	user, err := a.repo.FindByEmail(ctx, req.Email)
	if err == apperr.ErrUserNotFound {
		newUser := &models.User{
			Name:  req.Name,
			Email: req.Email,
			Image: req.Image,
		}
		createdUser, createErr := a.repo.Create(ctx, newUser)
		if createErr != nil {
			return nil, apperr.ErrInternalServer
		}
		return createdUser, nil
	}
	if err != nil {
		return nil, apperr.ErrInternalServer
	}
	return user, nil
}

func (a *AuthService) UpdateRole(ctx context.Context, role models.Role, email string) (*models.User, error) {
	return a.repo.UpdateRole(ctx, role, email) // ✅ returns *models.User now
}
