package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stripe/stripe-go/v78/webhook" // <-- ADDED WEBHOOK IMPORT
	"github.com/suryansh74/zomato/services/restaurant-service/internal/models"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/services"
	"github.com/suryansh74/zomato/services/shared/helper"
	"github.com/suryansh74/zomato/services/shared/middleware"
	"github.com/suryansh74/zomato/services/shared/token"
)

type OrderHandler struct {
	orderSrv *services.OrderService
}

func NewOrderHandler(orderSrv *services.OrderService) *OrderHandler {
	log.Println("Initializing OrderHandler")
	return &OrderHandler{orderSrv: orderSrv}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	userID := payload.User.ID

	var req models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	order, err := h.orderSrv.CreateOrder(r.Context(), userID, &req)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusCreated, map[string]any{
		"success": true,
		"message": "Order drafted successfully",
		"order":   order,
	})
}

func (h *OrderHandler) CreatePaymentSession(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	userID := payload.User.ID
	orderID := chi.URLParam(r, "id")

	checkoutURL, err := h.orderSrv.CreateStripeSession(r.Context(), orderID, userID)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"url": checkoutURL,
	})
}

func (h *OrderHandler) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Webhook Error reading body:", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// 🔒 SECURE WEBHOOK VERIFICATION ADDED HERE
	event, err := webhook.ConstructEventWithOptions(
		payload,
		r.Header.Get("Stripe-Signature"),
		"whsec_adaf4e2d8f3c94a9b346a681bf58a4ddf0dc2ef8b6b6416650ef2ff322ae3dfc",
		webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true},
	)
	if err != nil {
		log.Println("Webhook Error verifying signature:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if event.Type == "checkout.session.completed" {
		var session struct {
			Metadata map[string]string `json:"metadata"`
		}
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			log.Println("Error parsing session data:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		orderID := session.Metadata["order_id"]
		log.Println("🔔 Stripe Webhook Received: Payment Success for Order:", orderID)
		updatedOrder, err := h.orderSrv.ProcessPaymentSuccess(r.Context(), orderID)
		if err != nil {
			log.Println("Error marking order as paid:", err)
		} else {
			// ✅ TRIGGER THE REALTIME NOTIFICATION IN THE BACKGROUND
			go h.orderSrv.NotifyRestaurant(updatedOrder)
		}

		// TODO: This is where you will add RabbitMQ publishing later!
	}

	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	// Grab the logged-in user (who should be a restaurant owner)
	payload := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	println(payload.User.ID)
	// For now, we assume their User ID is tied to the Restaurant ID,
	// or they pass the Restaurant ID. Let's extract orderID and status.

	orderID := chi.URLParam(r, "id")

	var reqBody struct {
		Status       string `json:"status"`
		RestaurantID string `json:"restaurant_id"` // They will send this from the dashboard
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	err := h.orderSrv.UpdateOrderStatus(r.Context(), orderID, reqBody.RestaurantID, reqBody.Status)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// TODO LATER: Trigger Realtime Service to notify the CUSTOMER here!

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Order status updated to " + reqBody.Status,
	})
}

// Add to bottom of order_handler.go
func (h *OrderHandler) GetActiveOrders(w http.ResponseWriter, r *http.Request) {
	restaurantID := chi.URLParam(r, "id")

	orders, err := h.orderSrv.GetActiveOrders(r.Context(), restaurantID)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Ensure we return an empty array instead of null if there are no orders
	if orders == nil {
		orders = []models.Order{}
	}

	helper.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"orders": orders,
	})
}
