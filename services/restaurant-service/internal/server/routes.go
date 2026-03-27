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
	utilsClient := client.NewUtilsClient(s.cfg.UtilsServiceURL)

	// restaurant
	restaurantRepository := repositories.NewRestaurantRepository(s.client, s.cfg.DBName, s.cfg.CollectionName)
	restaurantService := services.NewRestaurantService(restaurantRepository)
	restaurantHandler := handlers.NewRestaurantHandler(restaurantService, utilsClient)

	// menu
	menuRepository := repositories.NewMenuRepository(s.client, s.cfg.DBName, "menu_items")
	menuService := services.NewMenuService(menuRepository)
	menuHandler := handlers.NewMenuHandler(menuService, restaurantService, utilsClient)

	// public routes
	s.router.Get("/api/restaurant/health", restaurantHandler.CheckHealth)
	s.router.Get("/api/restaurant/nearby", restaurantHandler.GetNearbyRestaurants)
	s.router.Get("/api/restaurant/{id}", restaurantHandler.GetRestaurantByID)

	// protected route for seller
	s.router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(s.tokenMaker))
		r.Use(restaurantMiddleware.IsRestaurantOwner())
		// restaurant routes
		r.Post("/api/restaurant/create", restaurantHandler.AddRestaurant)
		r.Get("/api/restaurant/read", restaurantHandler.GetRestaurant)
		r.Put("/api/restaurant/update", restaurantHandler.UpdateRestaurant)
		// menu routes
		r.Post("/api/menu", menuHandler.AddMenuItem)
		r.Get("/api/menu", menuHandler.GetMenuItems)
		r.Get("/api/menu/{id}", menuHandler.GetMenuItem)
		r.Put("/api/menu/{id}", menuHandler.UpdateMenuItem)
		r.Delete("/api/menu/{id}", menuHandler.DeleteMenuItem)
	})
}
