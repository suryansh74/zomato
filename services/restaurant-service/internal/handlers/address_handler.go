package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/models"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/services"
	"github.com/suryansh74/zomato/services/shared/helper"
	"github.com/suryansh74/zomato/services/shared/middleware"
	"github.com/suryansh74/zomato/services/shared/token"
)

type AddressHandler struct {
	addressSrv *services.AddressService
}

func NewAddressHandler(addressSrv *services.AddressService) *AddressHandler {
	log.Println("Initializing AddressHandler")
	return &AddressHandler{
		addressSrv: addressSrv,
	}
}

func (h *AddressHandler) AddAddress(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	userID := payload.User.ID

	var req models.AddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		helper.WriteJSON(w, http.StatusUnprocessableEntity, map[string]string{"errors": err.Error()})
		return
	}

	address, err := h.addressSrv.CreateAddress(r.Context(), userID, &req)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create address"})
		return
	}

	helper.WriteJSON(w, http.StatusCreated, map[string]any{
		"success": true,
		"message": "Address added successfully",
		"address": address,
	})
}

func (h *AddressHandler) GetMyAddresses(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	userID := payload.User.ID

	addresses, err := h.addressSrv.GetAddressesByUserID(r.Context(), userID)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch addresses"})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]any{
		"success":   true,
		"addresses": addresses,
	})
}

func (h *AddressHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	userID := payload.User.ID
	addressID := chi.URLParam(r, "id")

	err := h.addressSrv.DeleteAddress(r.Context(), addressID, userID)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete address"})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"success": "true",
		"message": "Address deleted successfully",
	})
}
