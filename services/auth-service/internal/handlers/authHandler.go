package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/suryansh74/zomato/services/auth-service/internal/models"
	services "github.com/suryansh74/zomato/services/auth-service/internal/serivces"
	"github.com/suryansh74/zomato/services/auth-service/internal/token"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var validate = validator.New()

type AuthHandler struct {
	tokenMaker          token.Maker
	accessTokenDuration time.Duration
	srv                 *services.AuthService
}

func NewAuthHandler(srv *services.AuthService, client *mongo.Client, dbName, collectionName string, tokenMaker token.Maker, accessTokenDuration time.Duration) *AuthHandler {
	return &AuthHandler{
		srv:                 srv,
		tokenMaker:          tokenMaker,
		accessTokenDuration: accessTokenDuration,
	}
}

// CheckHealth checks the health of the api
// ================================================================================
func (h *AuthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "auth-service is healthy",
	})
}

// Login user
// ================================================================================
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// getting incoming req body
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	// validate user
	if err := validate.Struct(req); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"errors": err.Error(),
		})
		return
	}

	user, err := h.srv.LoginOrCreate(&req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	token, err := h.tokenMaker.CreateToken(user, h.accessTokenDuration)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "login successfull",
		"token":   token,
		"user":    user,
	})
}

// helper functions
// ================================================================================

func writeJSON(w http.ResponseWriter, status int, data any) {
	responseBody := map[string]any{
		"data": data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(responseBody)
}
