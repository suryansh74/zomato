package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/suryansh74/zomato/services/auth-service/internal/helper"
	"github.com/suryansh74/zomato/services/auth-service/internal/middleware"
	"github.com/suryansh74/zomato/services/auth-service/internal/models"
	services "github.com/suryansh74/zomato/services/auth-service/internal/services"
	"github.com/suryansh74/zomato/services/auth-service/internal/token"
	"golang.org/x/oauth2"
)

var validate = validator.New()

type AuthHandler struct {
	tokenMaker          token.Maker
	accessTokenDuration time.Duration
	srv                 *services.AuthService
	oauthConfig         *oauth2.Config
}

func NewAuthHandler(srv *services.AuthService, tokenMaker token.Maker, accessTokenDuration time.Duration, oauthConfig *oauth2.Config) *AuthHandler {
	return &AuthHandler{
		srv:                 srv,
		tokenMaker:          tokenMaker,
		accessTokenDuration: accessTokenDuration,
		oauthConfig:         oauthConfig,
	}
}

func (h *AuthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "auth-service is healthy",
	})
}

// Login is the handler for the login endpoint
// ================================================================================
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// check if user wants login or signup
	// flow := r.URL.Query().Get("flow")
	// if flow == "" {
	// 	flow = "login"
	// }

	// send flow as state param
	url := h.oauthConfig.AuthCodeURL("login", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback is the handler for the google callback endpoint
// ===============================================================================
func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "missing code",
		})
		return
	}

	// get the flow back
	// flow := r.URL.Query().Get("state")

	token, err := h.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Println("token exchange error")
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	client := h.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to decode google response",
		})
		return
	}

	// create user if not exist
	user, err := h.srv.LoginOrCreate(r.Context(), &models.LoginRequest{
		Name:  googleUser.Name,
		Email: googleUser.Email,
		Image: googleUser.Picture,
	})
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	// create token
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
	// set cookie session
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
	// redirect
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
