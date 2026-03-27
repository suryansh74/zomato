package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	log.Println("Initializing RestaurantHandler")
	return &RestaurantHandler{
		srv:         srv,
		utilsClient: utilsClient,
	}
}

func (h *RestaurantHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler: CheckHealth called")

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "restaurant-service is healthy",
	})
}

// AddRestaurant adds a new restaurant
// ================================================================================
func (h *RestaurantHandler) AddRestaurant(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler: AddRestaurant called")

	email := r.Context().Value(middleware.UserContextKey).(*token.Payload).User.Email
	log.Println("AddRestaurant for email:", email)

	// check duplicate
	_, exists, err := h.srv.CheckIfOwnerHasRestaurant(r.Context(), email)
	if err != nil {
		log.Println("Error checking existing restaurant:", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if exists {
		log.Println("Restaurant already exists for email:", email)
		helper.WriteJSON(w, http.StatusConflict, map[string]string{"error": "you already have a restaurant"})
		return
	}

	// parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println("Error parsing multipart form:", err)
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid form data"})
		return
	}

	// get image file
	file, header, err := r.FormFile("image")
	if err != nil {
		log.Println("Image file missing:", err)
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "image is required"})
		return
	}
	defer file.Close()

	log.Println("Image received:", header.Filename)

	// forward cookie to utils service for auth
	cookie := r.Header.Get("Cookie")

	// call utils service to upload image
	imageURL, err := h.utilsClient.UploadImage(r.Context(), file, header.Filename, cookie)
	if err != nil {
		log.Println("Error uploading image:", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to upload image"})
		return
	}

	log.Println("Image uploaded successfully:", imageURL)

	// parse remaining form fields into struct
	req := models.RestaurantRequest{
		Name:             r.FormValue("name"),
		Description:      r.FormValue("description"),
		Image:            imageURL,
		Phone:            r.FormValue("phone"),
		FormattedAddress: r.FormValue("formatted_address"),
	}

	log.Println("Parsed form data:", req.Name, req.Phone)

	// parse coordinates
	lat, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
	if err != nil {
		log.Println("Invalid latitude:", err)
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid latitude"})
		return
	}
	lon, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
	if err != nil {
		log.Println("Invalid longitude:", err)
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid longitude"})
		return
	}
	req.Latitude = lat
	req.Longitude = lon

	log.Println("Coordinates parsed:", lat, lon)

	// validate
	if err := validate.Struct(req); err != nil {
		log.Println("Validation failed:", err)
		helper.WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": err.Error()})
		return
	}

	restaurant, err := h.srv.CreateRestaurant(r.Context(), email, &req)
	if err != nil {
		log.Println("Error creating restaurant:", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	log.Println("Restaurant created successfully:", restaurant.ID.Hex())
	helper.WriteJSON(w, http.StatusCreated, restaurant)
}

// GetRestaurant gets a single restaurant
// ================================================================================
func (h *RestaurantHandler) GetRestaurant(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler: GetRestaurant called")

	email := r.Context().Value(middleware.UserContextKey).(*token.Payload).User.Email
	log.Println("Fetching restaurant for email:", email)

	_, exists, err := h.srv.CheckIfOwnerHasRestaurant(r.Context(), email)
	if err != nil {
		log.Println("Error checking restaurant existence:", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if !exists {
		log.Println("Restaurant not found for email:", email)
		helper.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "restaurant not found"})
		return
	}

	restaurant, err := h.srv.GetRestaurant(r.Context(), email)
	if err != nil {
		log.Println("Error fetching restaurant:", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	log.Println("Restaurant fetched successfully:", restaurant.ID.Hex())
	helper.WriteJSON(w, http.StatusOK, restaurant)
}

func (h *RestaurantHandler) UpdateRestaurant(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler: UpdateRestaurant called")

	email := r.Context().Value(middleware.UserContextKey).(*token.Payload).User.Email

	var req models.UpdateRestaurantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Invalid request payload:", err)
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}

	log.Println("Update request received for email:", email)

	restaurant, err := h.srv.UpdateRestaurant(r.Context(), email, &req)
	if err != nil {
		log.Println("Error updating restaurant:", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	log.Println("Restaurant updated successfully:", restaurant.ID.Hex())
	helper.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Restaurant updated successfully",
		"data":    restaurant,
	})
}

// GetNearbyRestaurants fetches restaurants based on user location
func (h *RestaurantHandler) GetNearbyRestaurants(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler: GetNearbyRestaurants called")

	query := r.URL.Query()

	latStr := query.Get("latitude")
	lonStr := query.Get("longitude")

	if latStr == "" || lonStr == "" {
		log.Println("Missing latitude/longitude")
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "latitude and longitude are required"})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		log.Println("Invalid latitude format:", err)
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid latitude format"})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		log.Println("Invalid longitude format:", err)
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid longitude format"})
		return
	}

	radius := 5000.0
	if radStr := query.Get("radius"); radStr != "" {
		if parsedRad, err := strconv.ParseFloat(radStr, 64); err == nil {
			radius = parsedRad
		}
	}

	search := query.Get("search")

	var isOpenFilter *bool
	if openStr := query.Get("isOpen"); openStr != "" {
		if parsedOpen, err := strconv.ParseBool(openStr); err == nil {
			isOpenFilter = &parsedOpen
		}
	}

	log.Println("Nearby query params:", lat, lon, radius, search, isOpenFilter)

	restaurants, err := h.srv.GetNearbyRestaurants(r.Context(), lat, lon, radius, search, isOpenFilter)
	if err != nil {
		log.Println("Error fetching nearby restaurants:", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch nearby restaurants: " + err.Error()})
		return
	}

	log.Println("Nearby restaurants count:", len(restaurants))

	helper.WriteJSON(w, http.StatusOK, map[string]any{
		"count":       len(restaurants),
		"restaurants": restaurants,
	})
}

// GetRestaurantByID fetches a specific restaurant's details by its MongoDB ID
func (h *RestaurantHandler) GetRestaurantByID(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler: GetRestaurantByID called")

	id := chi.URLParam(r, "id")

	if id == "" {
		log.Println("Missing restaurant ID")
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "restaurant id is required"})
		return
	}

	log.Println("Fetching restaurant by ID:", id)

	restaurant, err := h.srv.GetRestaurantByID(r.Context(), id)
	if err != nil {
		log.Println("Error fetching restaurant by ID:", err)
		helper.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "restaurant not found or invalid ID"})
		return
	}

	log.Println("Restaurant fetched successfully:", restaurant.ID.Hex())

	helper.WriteJSON(w, http.StatusOK, map[string]any{
		"success":    true,
		"restaurant": restaurant,
	})
}
