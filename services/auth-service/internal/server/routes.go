package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/suryansh74/zomato/services/auth-service/internal/handlers"
	"github.com/suryansh74/zomato/services/auth-service/internal/middleware"
	"github.com/suryansh74/zomato/services/auth-service/internal/repositories"
	services "github.com/suryansh74/zomato/services/auth-service/internal/services"
)

func (s *Server) setupRoutes() {
	// handlers are created here and passed into routes
	authRepository := repositories.NewAuthRepository(s.client, s.cfg.DBName, s.cfg.CollectionName)
	authService := services.NewAuthService(authRepository)
	authHandler := handlers.NewAuthHandler(authService, s.tokenMaker, s.cfg.AccessTokenDuration)

	// health check
	s.router.Get("/api/health", authHandler.CheckHealth)
	s.router.Post("/api/auth/login", authHandler.Login)

	// protected route
	s.router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(s.tokenMaker))
		r.Get("/api/auth/profile", authHandler.Profile)
		r.Post("/api/auth/add_role", authHandler.AddRole)
	})
}
