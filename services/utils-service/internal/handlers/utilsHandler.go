package handlers

import (
	"image"
	"net/http"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/go-playground/validator/v10"
	"github.com/suryansh74/zomato/services/shared/helper"
	"github.com/suryansh74/zomato/services/utils-service/internal/models"
	services "github.com/suryansh74/zomato/services/utils-service/internal/services"
)

var validate = validator.New()

type UtilsHandler struct {
	srv *services.UtilsService
}

func NewUtilsHandler(srv *services.UtilsService) *UtilsHandler {
	return &UtilsHandler{
		srv: srv,
	}
}

func (h *UtilsHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "utils-service is healthy",
	})
}

func (h *UtilsHandler) ImageUpload(w http.ResponseWriter, r *http.Request) {
	// check method
	if r.Method != http.MethodPost {
		helper.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Only POST requests allowed",
		})
	}

	// Parse multipart form (max 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Get file from form (key = "image")
	file, _, err := r.FormFile("image")
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	defer file.Close()

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Get bounds
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	resp := &models.ImageDimensionResponse{
		Width:  width,
		Height: height,
	}

	helper.WriteJSON(w, http.StatusOK, map[string]any{
		"image_dimensions": resp,
	})
}
