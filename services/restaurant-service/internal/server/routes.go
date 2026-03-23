package server

import (
	"github.com/go-chi/chi/v5"
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
	restaurantHandler := handlers.NewRestaurantHandler(restaurantService)

	// health check
	s.router.Get("/api/restaurant/health", restaurantHandler.CheckHealth)

	// protected route
	s.router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(s.tokenMaker))
		r.Use(restaurantMiddleware.IsRestaurantOwner())
		// r.Post("/api/restaurant/addRestaurant", restaurantHandler.AddRestaurant)
	})
}
