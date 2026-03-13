package server

import (
	"github.com/suryansh74/zomato/services/auth-service/internal/handlers"
	"github.com/suryansh74/zomato/services/auth-service/internal/repositories"
	services "github.com/suryansh74/zomato/services/auth-service/internal/serivces"
)

func (s *Server) setupRoutes() {
	// handlers are created here and passed into routes
	authRepository := repositories.NewAuthRepository(s.client, s.cfg.DBName, s.cfg.CollectionName)
	authService := services.NewAuthService(authRepository)
	authHandler := handlers.NewAuthHandler(authService, s.client, s.cfg.DBName, s.cfg.CollectionName, s.tokenMaker, s.cfg.AccessTokenDuration)

	// health check
	s.router.Get("/api/health", authHandler.CheckHealth)
	s.router.Post("/api/auth/login", authHandler.Login)
}
