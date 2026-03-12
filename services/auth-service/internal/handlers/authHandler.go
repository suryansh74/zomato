package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/suryansh74/zomato/services/auth-service/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson" // ✅ v2, not v1
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var validate = validator.New()

type AuthHandler struct {
	client         *mongo.Client
	dbName         string
	collectionName string
}

func NewAuthHandler(client *mongo.Client, dbName, collectionName string) *AuthHandler {
	return &AuthHandler{
		client:         client,
		dbName:         dbName,
		collectionName: collectionName,
	}
}

func (h *AuthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "auth-service is healthy",
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// getting incoming req body
	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	if err := validate.Struct(req); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"errors": err.Error(),
		})
		return
	}

	// check email existance
	coll := h.client.Database(h.dbName).Collection(h.collectionName)

	var result bson.M
	err := coll.FindOne(r.Context(), bson.D{{Key: "email", Value: req.Email}}).Decode(&result)
	if err == mongo.ErrNoDocuments {

		// if user doesn't exist make record for it
		newUser := &models.User{
			Name:  req.Name,
			Email: req.Email,
			Image: req.Image,
			// Role:      "customer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		result, err := coll.InsertOne(context.TODO(), newUser)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "something went wrong",
			})
			return
		}

		fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
		writeJSON(w, http.StatusOK, map[string]any{
			"message": "login successful",
			"user":    result,
		})
		return
	}

	// user already existed
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "login successful",
		"user":    result,
	})

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "something went wrong",
		})
		return
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	responseBody := map[string]any{
		"data": data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(responseBody)
}
