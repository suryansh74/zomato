package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/client"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/handlers"
	restaurantMiddleware "github.com/suryansh74/zomato/services/restaurant-service/internal/middleware"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/repositories"
	services "github.com/suryansh74/zomato/services/restaurant-service/internal/services"
	"github.com/suryansh74/zomato/services/shared/middleware"
)

func (s *Server) setupRoutes() {
	// handlers are created here and passed into routes
	restaurantRepository := repositories.NewRestaurantRepository(s.client, s.cfg.DBName, s.cfg.CollectionName)
	restaurantService := services.NewRestaurantService(restaurantRepository)

	utilsClient := client.NewUtilsClient(s.cfg.UtilsServiceURL)
	restaurantHandler := handlers.NewRestaurantHandler(restaurantService, utilsClient)
	// health check
	s.router.Get("/api/restaurant/health", restaurantHandler.CheckHealth)

	// protected route
	s.router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(s.tokenMaker))
		r.Use(restaurantMiddleware.IsRestaurantOwner())
		r.Post("/api/restaurant/create", restaurantHandler.AddRestaurant)
		r.Get("/api/restaurant/read", restaurantHandler.GetRestaurant)
		r.Put("/api/restaurant/update", restaurantHandler.UpdateRestaurant)
	})
}
