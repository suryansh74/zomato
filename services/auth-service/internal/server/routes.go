package server

import (
	"github.com/suryansh74/zomato/services/auth-service/internal/handlers"
)

func (s *Server) setupRoutes() {
	// handlers are created here and passed into routes
	authHandler := handlers.NewAuthHandler(s.client, s.cfg.DBName, s.cfg.CollectionName)

	// health check
	s.router.Get("/api/health", authHandler.CheckHealth)
	s.router.Post("/api/auth/login", authHandler.Login)

	// // auth routes
	// s.router.Post("/api/auth/register", authHandler.Register)
	// s.router.Post("/api/auth/login", authHandler.Login)
}
