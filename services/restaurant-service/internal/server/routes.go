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

	// Repositories & Services & Handlers setup...
	restaurantRepository := repositories.NewRestaurantRepository(s.client, s.cfg.DBName, s.cfg.CollectionName)
	restaurantService := services.NewRestaurantService(restaurantRepository)
	restaurantHandler := handlers.NewRestaurantHandler(restaurantService, utilsClient)

	menuRepository := repositories.NewMenuRepository(s.client, s.cfg.DBName, "menu_items")
	menuService := services.NewMenuService(menuRepository)
	menuHandler := handlers.NewMenuHandler(menuService, restaurantService, utilsClient)

	// Health check can remain entirely public so your infrastructure (like AWS/Docker) can ping it
	s.router.Get("/api/restaurant/health", restaurantHandler.CheckHealth)

	// ==========================================
	// TIER 1: AUTHENTICATED ROUTES (Everyone)
	// ==========================================
	s.router.Group(func(r chi.Router) {
		// Gate 1: You MUST be logged in
		r.Use(middleware.AuthMiddleware(s.tokenMaker))

		// Customers (and owners acting as customers) can browse nearby food
		r.Get("/api/restaurant/nearby", restaurantHandler.GetNearbyRestaurants)
		r.Get("/api/restaurant/{id}", restaurantHandler.GetRestaurantByID)
		r.Get("/api/restaurant/{id}/menu", menuHandler.GetPublicMenu)

		// ==========================================
		// TIER 2: SELLER ROUTES (Owners Only)
		// ==========================================
		r.Group(func(ownerRouter chi.Router) {
			// Gate 2: You MUST be a restaurant owner
			ownerRouter.Use(restaurantMiddleware.IsRestaurantOwner())

			// Restaurant management
			ownerRouter.Post("/api/restaurant/create", restaurantHandler.AddRestaurant)
			ownerRouter.Get("/api/restaurant/read", restaurantHandler.GetRestaurant)
			ownerRouter.Put("/api/restaurant/update", restaurantHandler.UpdateRestaurant)

			// Menu management (The owner's specific dashboard endpoints)
			ownerRouter.Post("/api/menu", menuHandler.AddMenuItem)
			ownerRouter.Get("/api/menu", menuHandler.GetMenuItems)
			ownerRouter.Get("/api/menu/{id}", menuHandler.GetMenuItem)
			ownerRouter.Put("/api/menu/{id}", menuHandler.UpdateMenuItem)
			ownerRouter.Delete("/api/menu/{id}", menuHandler.DeleteMenuItem)
		})
	})
}
