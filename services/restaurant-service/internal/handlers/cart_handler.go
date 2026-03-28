package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/suryansh74/zomato/services/restaurant-service/apperr"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/models"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/services"
	"github.com/suryansh74/zomato/services/shared/helper"
	"github.com/suryansh74/zomato/services/shared/middleware"
	"github.com/suryansh74/zomato/services/shared/token"
)

type CartHandler struct {
	cartSrv       *services.CartService
	restaurantSrv *services.RestaurantService
	menuSrv       *services.MenuService
}

func NewCartHandler(cartSrv *services.CartService, restaurantSrv *services.RestaurantService, menuSrv *services.MenuService) *CartHandler {
	log.Println("Initializing CartHandler")
	return &CartHandler{
		cartSrv:       cartSrv,
		restaurantSrv: restaurantSrv,
		menuSrv:       menuSrv,
	}
}

// internal/handlers/cart_handler.go
func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	// 1. Get the User ID from the authenticated context
	payload := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	userID := payload.User.ID

	// 2. Decode the request body
	var req models.CartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	// 3. Set the UserID in the request struct
	req.UserID = userID

	if err := validate.Struct(req); err != nil {
		helper.WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": err.Error()})
		return
	}

	if req.ItemID == "" || req.RestaurantID == "" {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "item id and restaurant id are required"})
		return
	}

	// 4. Verify the restaurant exists (Optional, but good for data integrity)
	_, err := h.restaurantSrv.GetRestaurantByID(r.Context(), req.RestaurantID)
	if err != nil {
		helper.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "restaurant not found"})
		return
	}

	// 5. Verify the menu item exists and belongs to that restaurant
	_, err = h.menuSrv.GetMenuItemByID(r.Context(), req.ItemID, req.RestaurantID)
	if err != nil {
		helper.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "item not found for this restaurant"})
		return
	}

	// 6. FINALLY: Call the service to add it to the cart
	cartItem, err := h.cartSrv.AddToCart(r.Context(), &req)
	if err != nil {
		if err == apperr.ErrCartConflict {
			helper.WriteJSON(w, http.StatusConflict, map[string]string{
				"error": "You can order from only one restaurant at a time. Please clear your cart first.",
			})
			return
		}
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to add item to cart"})
		return
	}

	// 7. Return Success!
	helper.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Item added to cart",
		"cart":    cartItem,
	})
}

func (h *CartHandler) FetchCart(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler: FetchCart called")

	// 1. Get UserID from token
	payload := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	userID := payload.User.ID // Adjust this if you store ID differently in your token

	// 2. Fetch raw cart items
	cartItems, err := h.cartSrv.GetCartByUserID(r.Context(), userID)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch cart"})
		return
	}

	var populatedCart []map[string]any
	var subtotal float64 = 0
	var cartLength int = 0

	// 3. "Populate" the data manually
	for _, cart := range cartItems {
		// Fetch full details
		item, _ := h.menuSrv.GetMenuItemByID(r.Context(), cart.ItemID, cart.RestaurantID)
		restaurant, _ := h.restaurantSrv.GetRestaurantByID(r.Context(), cart.RestaurantID)

		// Create a populated object that perfectly matches the NodeJS output
		populatedItem := map[string]any{
			"_id":          cart.ID.Hex(),
			"userId":       cart.UserID,
			"restaurantId": restaurant, // Populated!
			"itemId":       item,       // Populated!
			"quantity":     cart.Quantity,
		}

		populatedCart = append(populatedCart, populatedItem)

		// Calculate Totals
		if item != nil {
			subtotal += item.Price * float64(cart.Quantity)
		}
		cartLength += cart.Quantity
	}

	// 4. Return exact JSON structure expected by frontend
	helper.WriteJSON(w, http.StatusOK, map[string]any{
		"success":    true,
		"cartLength": cartLength,
		"subtotal":   subtotal,
		"cart":       populatedCart,
	})
}

func (h *CartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	userID := payload.User.ID

	var req models.CartUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	if err := validate.Struct(req); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	err := h.cartSrv.UpdateQuantity(r.Context(), userID, req.ItemID, req.Action)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update cart"})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Cart updated"})
}

func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	userID := payload.User.ID

	err := h.cartSrv.ClearCart(r.Context(), userID)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to clear cart"})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Cart cleared"})
}
