package services

import (
	"context"
	"log"

	"github.com/suryansh74/zomato/services/auth-service/apperr"
	"github.com/suryansh74/zomato/services/auth-service/internal/models"
	"github.com/suryansh74/zomato/services/auth-service/internal/repositories"
)

type AuthService struct {
	repo repositories.AuthRepository
}

func NewAuthService(repo repositories.AuthRepository) *AuthService {
	log.Println("[AuthService] Initializing AuthService")
	return &AuthService{repo: repo}
}

func (a *AuthService) LoginOrCreate(ctx context.Context, req *models.LoginRequest) (*models.User, error) {
	log.Printf("[AuthService] LoginOrCreate called for email: %s\n", req.Email)

	user, err := a.repo.FindByEmail(ctx, req.Email)
	if err == apperr.ErrUserNotFound {
		log.Printf("[AuthService] User not found, creating new user for email: %s\n", req.Email)

		newUser := &models.User{
			Name:  req.Name,
			Email: req.Email,
			Image: req.Image,
		}

		createdUser, createErr := a.repo.Create(ctx, newUser)
		if createErr != nil {
			log.Printf("[AuthService] Error creating user for email %s: %v\n", req.Email, createErr)
			return nil, apperr.ErrInternalServer
		}

		log.Printf("[AuthService] User created successfully with email: %s\n", req.Email)
		return createdUser, nil
	}

	if err != nil {
		log.Printf("[AuthService] Error finding user for email %s: %v\n", req.Email, err)
		return nil, apperr.ErrInternalServer
	}

	log.Printf("[AuthService] User found for email: %s\n", req.Email)
	return user, nil
}

func (s *AuthService) UpdateRole(ctx context.Context, role models.Role, id string) (*models.User, error) {
	log.Printf("[AuthService] UpdateRole called for userID: %s with role: %s\n", id, role)

	user, err := s.repo.UpdateRole(ctx, role, id)
	if err != nil {
		log.Printf("[AuthService] Error updating role for userID %s: %v\n", id, err)
		return nil, err
	}

	log.Printf("[AuthService] Role updated successfully for userID: %s\n", id)
	return user, nil
}
