package main

import (
	"log"

	"github.com/suryansh74/zomato/services/restaurant-service/internal/config"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/db"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/server"
)

func main() {
	// 1. load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	// 2. connect to mongodb — client lives for entire app lifetime
	client := db.DBConnect(cfg.MongoURI)

	// 3. start server, pass client so handlers can use it later
	srv := server.NewServer(&cfg, client)
	srv.Start()
}
