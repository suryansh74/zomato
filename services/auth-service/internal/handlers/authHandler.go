package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/suryansh74/zomato/services/auth-service/internal/helper"
	"github.com/suryansh74/zomato/services/auth-service/internal/middleware"
	"github.com/suryansh74/zomato/services/auth-service/internal/models"
	services "github.com/suryansh74/zomato/services/auth-service/internal/services"
	"github.com/suryansh74/zomato/services/auth-service/internal/token"
)

var validate = validator.New()

type AuthHandler struct {
	tokenMaker          token.Maker
	accessTokenDuration time.Duration
	srv                 *services.AuthService
}

func NewAuthHandler(srv *services.AuthService, tokenMaker token.Maker, accessTokenDuration time.Duration) *AuthHandler {
	return &AuthHandler{
		srv:                 srv,
		tokenMaker:          tokenMaker,
		accessTokenDuration: accessTokenDuration,
	}
}

func (h *AuthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "auth-service is healthy",
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if err := validate.Struct(req); err != nil {
		helper.WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": err.Error()})
		return
	}

	user, err := h.srv.LoginOrCreate(r.Context(), &req)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// ✅ build TokenUser from User — includes role
	tokenUser := &models.TokenUser{
		Name:  user.Name,
		Email: user.Email,
		Image: user.Image,
		Role:  user.Role,
	}

	t, err := h.tokenMaker.CreateToken(tokenUser, h.accessTokenDuration)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    t,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.accessTokenDuration.Seconds()),
	})

	helper.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "login successful",
		"token":   t,
		"user":    user,
	})
}

func (h *AuthHandler) AddRole(w http.ResponseWriter, r *http.Request) {
	var role models.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if err := validate.Struct(role); err != nil {
		helper.WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": err.Error()})
		return
	}

	payload := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	email := payload.User.Email // ✅ payload.User not payload.LoginRequest

	updatedUser, err := h.srv.UpdateRole(r.Context(), role, email)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// ✅ regenerate token with updated role
	tokenUser := &models.TokenUser{
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
		Image: updatedUser.Image,
		Role:  updatedUser.Role,
	}

	t, err := h.tokenMaker.CreateToken(tokenUser, h.accessTokenDuration)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    t,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.accessTokenDuration.Seconds()),
	})

	helper.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "role updated successfully",
		"token":   t,
		"user":    updatedUser,
	})
}

func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	helper.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "profile successful",
		"payload": payload,
	})
}
