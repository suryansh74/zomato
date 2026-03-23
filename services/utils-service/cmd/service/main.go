package main

import (
	"log"

	"github.com/suryansh74/zomato/services/utils-service/internal/config"
	"github.com/suryansh74/zomato/services/utils-service/internal/server"
)

func main() {
	// developer chooses here 👇
	storageProvider := config.Cloudinary
	// storageProvider := config.GDrive

	cfg, err := config.LoadConfig(storageProvider)
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	srv := server.NewServer(&cfg)
	srv.Start()
}
