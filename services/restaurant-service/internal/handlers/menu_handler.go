package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/client"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/models"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/services"
	"github.com/suryansh74/zomato/services/shared/helper"
	"github.com/suryansh74/zomato/services/shared/middleware"
	"github.com/suryansh74/zomato/services/shared/token"
)

type MenuHandler struct {
	menuSrv       *services.MenuService
	restaurantSrv *services.RestaurantService // Needed to verify ownership
	utilsClient   *client.UtilsClient
}

func NewMenuHandler(menuSrv *services.MenuService, restaurantSrv *services.RestaurantService, utilsClient *client.UtilsClient) *MenuHandler {
	return &MenuHandler{
		menuSrv:       menuSrv,
		restaurantSrv: restaurantSrv,
		utilsClient:   utilsClient,
	}
}

// AddMenuItem parses form data, handles image upload, and saves the item
func (h *MenuHandler) AddMenuItem(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(middleware.UserContextKey).(*token.Payload).User.Email

	// 1. Verify the user has a restaurant and get the RestaurantID
	restaurantID, exists, err := h.restaurantSrv.CheckIfOwnerHasRestaurant(r.Context(), email)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if !exists {
		helper.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "you must create a restaurant first before adding menu items"})
		return
	}

	// 2. Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid form data"})
		return
	}

	var imageURL string

	// 3. Handle optional image file
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		cookie := r.Header.Get("Cookie")

		// Upload image
		imageURL, err = h.utilsClient.UploadImage(r.Context(), file, header.Filename, cookie)
		if err != nil {
			helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to upload image"})
			return
		}
	} else if err != http.ErrMissingFile {
		// If error is something other than "no file uploaded"
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "error parsing image file"})
		return
	}

	// 4. Parse strings to correct types for Price and IsAvailable
	priceStr := r.FormValue("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid price format"})
		return
	}

	isAvailableStr := r.FormValue("is_available")
	isAvailable, err := strconv.ParseBool(isAvailableStr)
	if err != nil {
		isAvailable = true // default to true if not provided or invalid
	}

	// 5. Populate Request Struct
	req := models.MenuItemRequest{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Image:       imageURL,
		Price:       price,
		IsAvailable: isAvailable,
	}

	// 6. Validate
	if err := validate.Struct(req); err != nil {
		helper.WriteJSON(w, http.StatusUnprocessableEntity, map[string]string{"errors": err.Error()})
		return
	}

	// 7. Save to DB
	menuItem, err := h.menuSrv.CreateMenuItem(r.Context(), restaurantID, &req)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusCreated, map[string]any{
		"message":   "Menu item added successfully",
		"menu_item": menuItem,
	})
}

// GetMenuItems fetches all items for the logged-in owner's restaurant
func (h *MenuHandler) GetMenuItems(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(middleware.UserContextKey).(*token.Payload).User.Email

	// Get RestaurantID
	restaurantID, exists, err := h.restaurantSrv.CheckIfOwnerHasRestaurant(r.Context(), email)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if !exists {
		helper.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "restaurant not found"})
		return
	}

	// Fetch Menu Items
	items, err := h.menuSrv.GetMenuItemsByRestaurant(r.Context(), restaurantID)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]any{
		"menu_items": items,
	})
}

func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(middleware.UserContextKey).(*token.Payload).User.Email
	id := chi.URLParam(r, "id") // Get ID from URL

	restaurantID, exists, err := h.restaurantSrv.CheckIfOwnerHasRestaurant(r.Context(), email)
	if err != nil || !exists {
		helper.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "unauthorized"})
		return
	}

	item, err := h.menuSrv.GetMenuItemByID(r.Context(), id, restaurantID)
	if err != nil {
		helper.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "menu item not found"})
		return
	}

	helper.WriteJSON(w, http.StatusOK, item) // Fixed double wrap
}

// UpdateMenuItem updates a single item by ID
func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(middleware.UserContextKey).(*token.Payload).User.Email
	id := chi.URLParam(r, "id")

	// 1. Verify Authorization
	restaurantID, exists, err := h.restaurantSrv.CheckIfOwnerHasRestaurant(r.Context(), email)
	if err != nil || !exists {
		helper.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "unauthorized"})
		return
	}

	// 2. Parse Multipart Form (Max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid form data"})
		return
	}

	var req models.UpdateMenuItemRequest

	// 3. Handle Optional Image Upload
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		cookie := r.Header.Get("Cookie")

		imageURL, uploadErr := h.utilsClient.UploadImage(r.Context(), file, header.Filename, cookie)
		if uploadErr != nil {
			helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to upload new image"})
			return
		}
		req.Image = &imageURL
	} else if err != http.ErrMissingFile {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "error parsing image file"})
		return
	}

	// 4. Handle Text/Number Fields safely
	// We check if the field was actually sent in the request before assigning the pointer
	if r.MultipartForm != nil && r.MultipartForm.Value != nil {
		if vals, ok := r.MultipartForm.Value["name"]; ok && len(vals) > 0 {
			name := vals[0]
			req.Name = &name
		}

		if vals, ok := r.MultipartForm.Value["description"]; ok && len(vals) > 0 {
			desc := vals[0]
			req.Description = &desc
		}

		if vals, ok := r.MultipartForm.Value["price"]; ok && len(vals) > 0 {
			if price, err := strconv.ParseFloat(vals[0], 64); err == nil {
				req.Price = &price
			}
		}

		if vals, ok := r.MultipartForm.Value["is_available"]; ok && len(vals) > 0 {
			if isAvailable, err := strconv.ParseBool(vals[0]); err == nil {
				req.IsAvailable = &isAvailable
			}
		}
	}

	// 5. Update Database
	item, err := h.menuSrv.UpdateMenuItem(r.Context(), id, restaurantID, &req)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusOK, item)
}

// DeleteMenuItem deletes a single item by ID
func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(middleware.UserContextKey).(*token.Payload).User.Email
	id := chi.URLParam(r, "id") // Get ID from URL

	restaurantID, exists, err := h.restaurantSrv.CheckIfOwnerHasRestaurant(r.Context(), email)
	if err != nil || !exists {
		helper.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "unauthorized"})
		return
	}

	err = h.menuSrv.DeleteMenuItem(r.Context(), id, restaurantID)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete item"})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Menu item deleted successfully"})
}
