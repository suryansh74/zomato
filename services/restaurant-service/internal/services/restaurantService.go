package services

import (
	"context"

	"github.com/suryansh74/zomato/services/restaurant-service/internal/models"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/repositories"
)

type RestaurantService struct {
	repo repositories.RestaurantRepository
}

func NewRestaurantService(repo repositories.RestaurantRepository) *RestaurantService {
	return &RestaurantService{repo: repo}
}

func (s *RestaurantService) CheckIfOwnerHasRestaurant(ctx context.Context, email string) (bool, error) {
	return s.repo.CheckIfOwnerHasRestaurant(ctx, email)
}

func (s *RestaurantService) CreateRestaurant(ctx context.Context, ownerEmail string, req *models.RestaurantRequest) (*models.Restaurant, error) {
	restaurant := &models.Restaurant{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		OwnerEmail:  ownerEmail,
		Phone:       req.Phone,
		IsVerified:  false,
		IsOpen:      false,
		AutoLocation: models.GeoJSONPoint{
			Type:             "Point",
			Coordinates:      []float64{req.Longitude, req.Latitude}, // longitude first
			FormattedAddress: req.FormattedAddress,
		},
	}
	return s.repo.CreateRestaurant(ctx, restaurant)
}
