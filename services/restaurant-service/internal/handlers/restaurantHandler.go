package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/client"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/models"
	services "github.com/suryansh74/zomato/services/restaurant-service/internal/services"
	"github.com/suryansh74/zomato/services/shared/helper"
	"github.com/suryansh74/zomato/services/shared/middleware"
	"github.com/suryansh74/zomato/services/shared/token"
)

var validate = validator.New()

type RestaurantHandler struct {
	srv         *services.RestaurantService
	utilsClient *client.UtilsClient
}

func NewRestaurantHandler(srv *services.RestaurantService, utilsClient *client.UtilsClient) *RestaurantHandler {
	return &RestaurantHandler{
		srv:         srv,
		utilsClient: utilsClient,
	}
}

func (h *RestaurantHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "restaurant-service is healthy",
	})
}

func (h *RestaurantHandler) AddRestaurant(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(middleware.UserContextKey).(*token.Payload).User.Email

	// check duplicate
	exists, err := h.srv.CheckIfOwnerHasRestaurant(r.Context(), email)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if exists {
		helper.WriteJSON(w, http.StatusConflict, map[string]string{"error": "you already have a restaurant"})
		return
	}

	// parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid form data"})
		return
	}

	// get image file
	file, header, err := r.FormFile("image")
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "image is required"})
		return
	}
	defer file.Close()

	// forward cookie to utils service for auth
	cookie := r.Header.Get("Cookie")

	// call utils service to upload image
	imageURL, err := h.utilsClient.UploadImage(r.Context(), file, header.Filename, cookie)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to upload image"})
		return
	}

	// parse remaining form fields into struct
	req := models.RestaurantRequest{
		Name:             r.FormValue("name"),
		Description:      r.FormValue("description"),
		Image:            imageURL, // from utils service
		Phone:            r.FormValue("phone"),
		FormattedAddress: r.FormValue("formatted_address"),
	}

	// parse coordinates
	lat, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid latitude"})
		return
	}
	lon, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid longitude"})
		return
	}
	req.Latitude = lat
	req.Longitude = lon

	// validate
	if err := validate.Struct(req); err != nil {
		helper.WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": err.Error()})
		return
	}

	restaurant, err := h.srv.CreateRestaurant(r.Context(), email, &req)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusCreated, restaurant)
}
