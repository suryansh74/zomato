package handlers

import (
	"context"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-playground/validator/v10"
	"github.com/suryansh74/zomato/services/shared/helper"
	services "github.com/suryansh74/zomato/services/utils-service/internal/services"
)

var validate = validator.New()

type UtilsHandler struct {
	srv *services.UtilsService
	cld *cloudinary.Cloudinary
	ctx context.Context
}

func NewUtilsHandler(srv *services.UtilsService, cld *cloudinary.Cloudinary, ctx context.Context) *UtilsHandler {
	return &UtilsHandler{
		srv: srv,
		cld: cld,
		ctx: ctx,
	}
}

func (h *UtilsHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "utils-service is healthy",
	})
}

// ImageUpload handles the HTTP request and passes data to the service layer
func (h *UtilsHandler) ImageUpload(w http.ResponseWriter, r *http.Request) {
	// Protect the HTTP server from excessively large payloads
	r.Body = http.MaxBytesReader(w, r.Body, services.MaxUploadSize+1024)

	if err := r.ParseMultipartForm(services.MaxUploadSize); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "File upload exceeds limits or is malformed",
		})
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Failed to get image from request. Ensure form field is named 'image'",
		})
		return
	}
	defer file.Close()

	// Pass file to the service layer for logic and validation
	secureURL, err := h.srv.ProcessAndUploadImage(r.Context(), h.cld, file, fileHeader)
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Return the final URL
	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message":   "Image uploaded successfully",
		"image_url": secureURL,
	})
}
