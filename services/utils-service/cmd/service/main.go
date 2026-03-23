package main

import (
	"log"

	"github.com/suryansh74/zomato/services/utils-service/internal/config"
	"github.com/suryansh74/zomato/services/utils-service/internal/server"
)

func main() {
	// 1. load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	// 2. start server, pass client so handlers can use it later
	srv := server.NewServer(&cfg)
	srv.Start()
}
