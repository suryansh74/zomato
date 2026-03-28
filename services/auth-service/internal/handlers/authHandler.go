package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	authModel "github.com/suryansh74/zomato/services/auth-service/internal/models"
	services "github.com/suryansh74/zomato/services/auth-service/internal/services"
	"github.com/suryansh74/zomato/services/shared/helper"
	"github.com/suryansh74/zomato/services/shared/middleware"
	tokenModel "github.com/suryansh74/zomato/services/shared/models"
	"github.com/suryansh74/zomato/services/shared/token"
	"golang.org/x/oauth2"
)

var validate = validator.New()

type AuthHandler struct {
	tokenMaker          token.Maker
	accessTokenDuration time.Duration
	srv                 *services.AuthService
	oauthConfig         *oauth2.Config
	isDev               bool
	frontendURL         string
}

func NewAuthHandler(srv *services.AuthService, tokenMaker token.Maker, accessTokenDuration time.Duration, oauthConfig *oauth2.Config, isDev bool, frontendURL string) *AuthHandler {
	log.Println("[AuthHandler] Initializing AuthHandler")
	return &AuthHandler{
		srv:                 srv,
		tokenMaker:          tokenMaker,
		accessTokenDuration: accessTokenDuration,
		oauthConfig:         oauthConfig,
		isDev:               isDev,
		frontendURL:         frontendURL,
	}
}

func (h *AuthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	log.Println("[AuthHandler] Health check endpoint hit")
	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "auth-service is healthy",
	})
}

// Login handler
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("[AuthHandler] Login endpoint hit, redirecting to Google OAuth")

	url := h.oauthConfig.AuthCodeURL("login", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback handler
func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	log.Println("[AuthHandler] GoogleCallback triggered")

	code := r.URL.Query().Get("code")
	if code == "" {
		log.Println("[AuthHandler] Missing code in callback")
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "missing code",
		})
		return
	}

	log.Println("[AuthHandler] Exchanging code for token")
	token, err := h.oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Printf("[AuthHandler] Token exchange error: %v\n", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	client := h.oauthConfig.Client(r.Context(), token)

	log.Println("[AuthHandler] Fetching user info from Google")
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("[AuthHandler] Error fetching user info: %v\n", err)
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
		log.Printf("[AuthHandler] Failed to decode Google response: %v\n", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to decode google response",
		})
		return
	}

	log.Printf("[AuthHandler] Google user fetched: email=%s\n", googleUser.Email)

	// create or login user
	user, err := h.srv.LoginOrCreate(r.Context(), &authModel.LoginRequest{
		Name:  googleUser.Name,
		Email: googleUser.Email,
		Image: googleUser.Picture,
	})
	if err != nil {
		log.Printf("[AuthHandler] LoginOrCreate failed: %v\n", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	log.Printf("[AuthHandler] User authenticated: id=%s email=%s\n", user.ID.Hex(), user.Email)

	// create token
	tokenUser := &tokenModel.TokenUser{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
		Image: user.Image,
		Role:  user.Role,
	}

	log.Println("[AuthHandler] Creating access token")
	t, err := h.tokenMaker.CreateToken(tokenUser, h.accessTokenDuration)
	if err != nil {
		log.Printf("[AuthHandler] Token creation failed: %v\n", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	log.Println("[AuthHandler] Setting session cookie")
	h.setSessionCookie(w, t)

	if user.Role == "" {
		log.Println("[AuthHandler] New user detected, redirecting to role selection")
		http.Redirect(
			w,
			r,
			fmt.Sprintf("%s/select-role?fresh=true", h.frontendURL),
			http.StatusTemporaryRedirect,
		)
	} else {
		log.Println("[AuthHandler] Existing user, redirecting to home")
		http.Redirect(
			w,
			r,
			fmt.Sprintf("%s/?fresh=true", h.frontendURL),
			http.StatusTemporaryRedirect,
		)
	}
}

func (h *AuthHandler) AddRole(w http.ResponseWriter, r *http.Request) {
	log.Println("[AuthHandler] AddRole endpoint hit")

	var role authModel.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		log.Println("[AuthHandler] Invalid request body for AddRole")
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(role); err != nil {
		log.Printf("[AuthHandler] Validation failed: %v\n", err)
		helper.WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": err.Error()})
		return
	}

	payload := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	userID := payload.User.ID

	log.Printf("[AuthHandler] Updating role for userID: %s\n", userID)

	updatedUser, err := h.srv.UpdateRole(r.Context(), role, userID)
	if err != nil {
		log.Printf("[AuthHandler] UpdateRole failed: %v\n", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	log.Printf("[AuthHandler] Role updated successfully for userID: %s\n", userID)

	tokenUser := &tokenModel.TokenUser{
		ID:    updatedUser.ID.Hex(),
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
		Image: updatedUser.Image,
		Role:  updatedUser.Role,
	}

	log.Println("[AuthHandler] Regenerating token after role update")
	t, err := h.tokenMaker.CreateToken(tokenUser, h.accessTokenDuration)
	if err != nil {
		log.Printf("[AuthHandler] Token creation failed after role update: %v\n", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	log.Println("[AuthHandler] Setting session cookie after role update")
	h.setSessionCookie(w, t)

	helper.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "role updated successfully",
		"user":    updatedUser,
	})
}

func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	log.Println("[AuthHandler] Profile endpoint hit")

	payload := r.Context().Value(middleware.UserContextKey).(*token.Payload)

	log.Printf("[AuthHandler] Returning profile for userID: %s\n", payload.User.ID)

	helper.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "profile successful",
		"payload": payload,
	})
}

func (h *AuthHandler) setSessionCookie(w http.ResponseWriter, token string) {
	log.Println("[AuthHandler] Setting session cookie")

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   !h.isDev,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.accessTokenDuration.Seconds()),
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	log.Println("[AuthHandler] Logout endpoint hit")

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   !h.isDev,
		MaxAge:   -1,
	})

	log.Println("[AuthHandler] Session cookie cleared")

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "logged out",
	})
}
