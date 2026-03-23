package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	services "github.com/suryansh74/zomato/services/restaurant-service/internal/services"
	"github.com/suryansh74/zomato/services/shared/helper"
)

var validate = validator.New()

type RestaurantHandler struct {
	srv *services.RestaurantService
}

func NewRestaurantHandler(srv *services.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{
		srv: srv,
	}
}

func (h *RestaurantHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "restaurant-service is healthy",
	})
}
