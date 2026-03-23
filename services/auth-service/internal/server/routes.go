package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/suryansh74/zomato/services/auth-service/internal/handlers"
	"github.com/suryansh74/zomato/services/auth-service/internal/repositories"
	services "github.com/suryansh74/zomato/services/auth-service/internal/services"
	"github.com/suryansh74/zomato/services/shared/middleware"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func (s *Server) setupRoutes() {
	oauthConfig := &oauth2.Config{
		ClientID:     s.cfg.GoogleClientID,
		ClientSecret: s.cfg.GoogleClientSecret,
		RedirectURL:  s.cfg.GoogleRedirectURL,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	// handlers are created here and passed into routes
	authRepository := repositories.NewAuthRepository(s.client, s.cfg.DBName, s.cfg.CollectionName)
	authService := services.NewAuthService(authRepository)
	authHandler := handlers.NewAuthHandler(authService, s.tokenMaker, s.cfg.AccessTokenDuration, oauthConfig, s.cfg.IsDev, s.cfg.FrontendURL)

	// health check
	s.router.Get("/api/auth/health", authHandler.CheckHealth)
	s.router.Get("/api/auth/login", authHandler.Login)
	s.router.Get("/api/auth/google/callback", authHandler.GoogleCallback)

	// protected route
	s.router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(s.tokenMaker))
		r.Get("/api/auth/profile", authHandler.Profile)
		r.Post("/api/auth/add_role", authHandler.AddRole)
		r.Post("/api/auth/logout", authHandler.Logout)
	})
}
