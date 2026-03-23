package server

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/suryansh74/zomato/services/shared/middleware"
	"github.com/suryansh74/zomato/services/utils-service/internal/client"
	"github.com/suryansh74/zomato/services/utils-service/internal/handlers"
	services "github.com/suryansh74/zomato/services/utils-service/internal/services"
)

func (s *Server) setupRoutes() {
	// handlers are created here and passed into routes
	// utilsService := services.NewUtilsService()
	// utilsHandler := handlers.NewUtilsHandler(utilsService)
	utilsService := services.NewUtilsService()
	cld, ctx, err := client.NewCloudinary(s.cfg.CloudinaryURL)
	if err != nil {
		log.Fatal("cloudinary init failed:", err)
	}

	utilsHandler := handlers.NewUtilsHandler(utilsService, cld, ctx)

	// health check
	s.router.Get("/api/utils/health", utilsHandler.CheckHealth)

	// protected route
	s.router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(s.tokenMaker))
		r.Post("/api/utils/upload", utilsHandler.ImageUpload)
	})
}
