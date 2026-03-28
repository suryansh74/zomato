package handlers

import (
	"log"
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
	restaurantSrv *services.RestaurantService
	utilsClient   *client.UtilsClient
}

func NewMenuHandler(menuSrv *services.MenuService, restaurantSrv *services.RestaurantService, utilsClient *client.UtilsClient) *MenuHandler {
	log.Println("Initializing MenuHandler")
	return &MenuHandler{
		menuSrv:       menuSrv,
		restaurantSrv: restaurantSrv,
		utilsClient:   utilsClient,
	}
}

// AddMenuItem parses form data, handles image upload, and saves the item
func (h *MenuHandler) AddMenuItem(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler: AddMenuItem called")

	email := r.Context().Value(middleware.UserContextKey).(*token.Payload).User.Email
	log.Println("AddMenuItem for email:", email)

	restaurantID, exists, err := h.restaurantSrv.CheckIfOwnerHasRestaurant(r.Context(), email)
	if err != nil {
		log.Println("Error checking restaurant ownership:", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if !exists {
		log.Println("No restaurant found for email:", email)
		helper.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "you must create a restaurant first before adding menu items"})
		return
	}

	log.Println("Restaurant ID:", restaurantID)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println("Error parsing multipart form:", err)
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid form data"})
		return
	}

	var imageURL string

	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		log.Println("Image received:", header.Filename)

		cookie := r.Header.Get("Cookie")

		imageURL, err = h.utilsClient.UploadImage(r.Context(), file, header.Filename, cookie)
		if err != nil {
			log.Println("Error uploading image:", err)
			helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to upload image"})
			return
		}
		log.Println("Image uploaded:", imageURL)

	} else if err != http.ErrMissingFile {
		log.Println("Error parsing image file:", err)
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "error parsing image file"})
		return
	} else {
		log.Println("No image provided")
	}

	priceStr := r.FormValue("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Println("Invalid price:", priceStr)
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid price format"})
		return
	}

	isAvailableStr := r.FormValue("is_available")
	isAvailable, err := strconv.ParseBool(isAvailableStr)
	if err != nil {
		log.Println("Invalid is_available, defaulting to true:", isAvailableStr)
		isAvailable = true
	}

	req := models.MenuItemRequest{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Image:       imageURL,
		Price:       price,
		IsAvailable: isAvailable,
	}

	log.Println("Parsed menu item:", req.Name, req.Price)

	if err := validate.Struct(req); err != nil {
		log.Println("Validation failed:", err)
		helper.WriteJSON(w, http.StatusUnprocessableEntity, map[string]string{"errors": err.Error()})
		return
	}

	menuItem, err := h.menuSrv.CreateMenuItem(r.Context(), restaurantID, &req)
	if err != nil {
		log.Println("Error creating menu item:", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	log.Println("Menu item created:", menuItem.ID.Hex())

	helper.WriteJSON(w, http.StatusCreated, map[string]any{
		"message":   "Menu item added successfully",
		"menu_item": menuItem,
	})
}

// GetMenuItems fetches all items for the logged-in owner's restaurant
func (h *MenuHandler) GetMenuItems(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler: GetMenuItems called")

	email := r.Context().Value(middleware.UserContextKey).(*token.Payload).User.Email
	log.Println("Fetching menu items for email:", email)

	restaurantID, exists, err := h.restaurantSrv.CheckIfOwnerHasRestaurant(r.Context(), email)
	if err != nil {
		log.Println("Error checking restaurant:", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if !exists {
		log.Println("Restaurant not found for email:", email)
		helper.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "restaurant not found"})
		return
	}

	items, err := h.menuSrv.GetMenuItemsByRestaurant(r.Context(), restaurantID)
	if err != nil {
		log.Println("Error fetching menu items:", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	log.Println("Menu items count:", len(items))

	helper.WriteJSON(w, http.StatusOK, map[string]any{
		"menu_items": items,
	})
}

func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler: GetMenuItem called")

	email := r.Context().Value(middleware.UserContextKey).(*token.Payload).User.Email
	id := chi.URLParam(r, "id")

	log.Println("Fetching menu item id:", id, "for email:", email)

	restaurantID, exists, err := h.restaurantSrv.CheckIfOwnerHasRestaurant(r.Context(), email)
	if err != nil || !exists {
		log.Println("Unauthorized access for menu item:", id)
		helper.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "unauthorized"})
		return
	}

	item, err := h.menuSrv.GetMenuItemByID(r.Context(), id, restaurantID)
	if err != nil {
		log.Println("Menu item not found:", id)
		helper.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "menu item not found"})
		return
	}

	log.Println("Menu item fetched:", item.ID.Hex())
	helper.WriteJSON(w, http.StatusOK, item)
}

// UpdateMenuItem updates a single item by ID
func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler: UpdateMenuItem called")

	email := r.Context().Value(middleware.UserContextKey).(*token.Payload).User.Email
	id := chi.URLParam(r, "id")

	log.Println("Updating menu item id:", id, "for email:", email)

	restaurantID, exists, err := h.restaurantSrv.CheckIfOwnerHasRestaurant(r.Context(), email)
	if err != nil || !exists {
		log.Println("Unauthorized update attempt:", id)
		helper.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "unauthorized"})
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println("Error parsing form:", err)
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid form data"})
		return
	}

	var req models.UpdateMenuItemRequest

	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		log.Println("New image received:", header.Filename)

		cookie := r.Header.Get("Cookie")

		imageURL, uploadErr := h.utilsClient.UploadImage(r.Context(), file, header.Filename, cookie)
		if uploadErr != nil {
			log.Println("Error uploading new image:", uploadErr)
			helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to upload new image"})
			return
		}
		req.Image = &imageURL
	} else if err != http.ErrMissingFile {
		log.Println("Error parsing image file:", err)
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "error parsing image file"})
		return
	}

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

	item, err := h.menuSrv.UpdateMenuItem(r.Context(), id, restaurantID, &req)
	if err != nil {
		log.Println("Error updating menu item:", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	log.Println("Menu item updated:", item.ID.Hex())
	helper.WriteJSON(w, http.StatusOK, item)
}

// DeleteMenuItem deletes a single item by ID
func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler: DeleteMenuItem called")

	email := r.Context().Value(middleware.UserContextKey).(*token.Payload).User.Email
	id := chi.URLParam(r, "id")

	log.Println("Deleting menu item id:", id, "for email:", email)

	restaurantID, exists, err := h.restaurantSrv.CheckIfOwnerHasRestaurant(r.Context(), email)
	if err != nil || !exists {
		log.Println("Unauthorized delete attempt:", id)
		helper.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "unauthorized"})
		return
	}

	err = h.menuSrv.DeleteMenuItem(r.Context(), id, restaurantID)
	if err != nil {
		log.Println("Error deleting menu item:", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete item"})
		return
	}

	log.Println("Menu item deleted:", id)
	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Menu item deleted successfully"})
}

// GetPublicMenu fetches all available menu items for a specific restaurant ID (For Customers)
func (h *MenuHandler) GetPublicMenu(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler: GetPublicMenu called")

	restaurantID := chi.URLParam(r, "id")
	if restaurantID == "" {
		log.Println("Missing restaurant ID")
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "restaurant id is required"})
		return
	}

	log.Println("Fetching public menu for restaurantID:", restaurantID)

	items, err := h.menuSrv.GetMenuItemsByRestaurant(r.Context(), restaurantID)
	if err != nil {
		log.Println("Error fetching menu items:", err)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch menu items"})
		return
	}

	var availableItems []models.MenuItem
	for _, item := range items {
		if item.IsAvailable {
			availableItems = append(availableItems, item)
		}
	}

	if availableItems == nil {
		availableItems = []models.MenuItem{}
	}

	log.Println("Available menu items count:", len(availableItems))

	helper.WriteJSON(w, http.StatusOK, map[string]any{
		"success":    true,
		"menu_items": availableItems,
	})
}
