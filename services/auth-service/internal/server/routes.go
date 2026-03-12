package server

import (
	"github.com/suryansh74/zomato/services/auth-service/internal/handlers"
	"github.com/suryansh74/zomato/services/auth-service/internal/repositories"
	"github.com/suryansh74/zomato/services/auth-service/internal/serivces"
)

func (s *Server) setupRoutes() {
	// handlers are created here and passed into routes
	authRepository := repositories.NewAuthRepository(s.client, s.cfg.DBName, s.cfg.CollectionName)
	authService := serivces.NewAuthService(authRepository)
	authHandler := handlers.NewAuthHandler(authService, s.client, s.cfg.DBName, s.cfg.CollectionName)

	// health check
	s.router.Get("/api/health", authHandler.CheckHealth)
	s.router.Post("/api/auth/login", authHandler.Login)
}
