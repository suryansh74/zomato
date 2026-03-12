package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/suryansh74/zomato/services/auth-service/internal/config"
	"github.com/suryansh74/zomato/services/auth-service/internal/token"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Server struct {
	cfg        *config.Config
	router     *chi.Mux
	client     *mongo.Client
	tokenMaker token.Maker
}

func NewServer(cfg *config.Config, client *mongo.Client) *Server {
	tokenMaker, err := token.NewPasetoMaker(cfg.TokenSymmetricKey)
	if err != nil {
		panic(err)
	}
	return &Server{
		cfg:        cfg,
		router:     chi.NewRouter(),
		client:     client,
		tokenMaker: tokenMaker,
	}
}

func (s *Server) Start() {
	// 1. attach global middleware
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	// 2. register routes
	s.setupRoutes()

	// 3. start server
	addr := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)
	fmt.Println("auth-service running on", addr)
	http.ListenAndServe(addr, s.router)
}
