package services

import (
	"context"

	"github.com/suryansh74/zomato/services/restaurant-service/internal/models"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/repositories"
)

type RestaurantService struct {
	repo repositories.RestaurantRepository
}

// NewRestaurantService constructor for restaurantService
// ================================================================================
func NewRestaurantService(repo repositories.RestaurantRepository) *RestaurantService {
	return &RestaurantService{repo: repo}
}

// CheckIfOwnerHasRestaurant
// ================================================================================
func (s *RestaurantService) CheckIfOwnerHasRestaurant(ctx context.Context, email string) (bool, error) {
	return s.repo.CheckIfOwnerHasRestaurant(ctx, email)
}

// CreateRestaurant
// ================================================================================
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

// GetRestaurant
// ================================================================================
func (s *RestaurantService) GetRestaurant(ctx context.Context, email string) (*models.Restaurant, error) {
	return s.repo.GetRestaurant(ctx, email)
}

// UpdateRestaurant
// ===============================================================================
func (s *RestaurantService) UpdateRestaurant(ctx context.Context, email string, req *models.UpdateRestaurantRequest) (*models.Restaurant, error) {
	return s.repo.UpdateRestaurant(ctx, email, req)
}
