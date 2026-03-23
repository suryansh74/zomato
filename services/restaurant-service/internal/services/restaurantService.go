package services

import (
	"github.com/suryansh74/zomato/services/restaurant-service/internal/repositories"
)

type RestaurantService struct {
	repo repositories.RestaurantRepository
}

func NewRestaurantService(repo repositories.RestaurantRepository) *RestaurantService {
	return &RestaurantService{repo: repo}
}
