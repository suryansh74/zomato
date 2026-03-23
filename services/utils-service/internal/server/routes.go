package server

import (
	"context"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/suryansh74/zomato/services/shared/middleware"
	"github.com/suryansh74/zomato/services/utils-service/internal/adapters"
	"github.com/suryansh74/zomato/services/utils-service/internal/client"
	"github.com/suryansh74/zomato/services/utils-service/internal/handlers"
	services "github.com/suryansh74/zomato/services/utils-service/internal/services"
)

func (s *Server) setupRoutes() {
	var storageAdapter adapters.StorageAdapter

	// Example: Choose adapter based on config
	if s.cfg.StorageProvider == "cloudinary" {
		cld, _, err := client.NewCloudinary(s.cfg.CloudinaryURL)
		if err != nil {
			log.Fatal(err)
		}
		storageAdapter = adapters.NewCloudinaryAdapter(cld)
		log.Println("Using Cloudinary for storage")

	} else if s.cfg.StorageProvider == "gdrive" {
		// --- ADD THESE DEBUG LINES ---
		log.Printf("DEBUG: Using Credentials Path: '%s'\n", s.cfg.GoogleCredentialsPath)
		log.Printf("DEBUG: Using Folder ID: '%s'\n", s.cfg.GoogleDriveFolderID)
		if s.cfg.GoogleDriveFolderID == "" {
			log.Fatal("CRITICAL: Folder ID is empty! Check your .env file parsing.")
		}
		// ------------------------------

		gdrive, err := adapters.NewGoogleDriveAdapter(context.Background(), s.cfg.GoogleCredentialsPath, s.cfg.GoogleDriveFolderID)
		if err != nil {
			log.Fatal(err)
		}
		storageAdapter = gdrive
		log.Println("Using Google Drive for storage")
	}

	// Inject the chosen adapter into the service
	utilsService := services.NewUtilsService(storageAdapter)

	// Note: You no longer need to pass cld and ctx to the handler!
	utilsHandler := handlers.NewUtilsHandler(utilsService)
	// health check
	s.router.Get("/api/utils/health", utilsHandler.CheckHealth)

	// protected route
	s.router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(s.tokenMaker))
		r.Post("/api/utils/upload", utilsHandler.ImageUpload)
	})
}
