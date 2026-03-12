// Package token response for token functionality
package token

import (
	"errors"
	"time"

	"github.com/suryansh74/zomato/services/auth-service/internal/models"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token is expired")
)

// Payload contain payload data of token
type Payload struct {
	LoginRequest *models.LoginRequest `json:"user"`
	IssuedAt     time.Time            `json:"issued_at"`
	ExpiredAt    time.Time            `json:"expired_at"`
}

func NewPayload(LoginRequest *models.LoginRequest, duration time.Duration) (*Payload, error) {
	return &Payload{
		LoginRequest: LoginRequest,
		IssuedAt:     time.Now(),
		ExpiredAt:    time.Now().Add(duration),
	}, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
