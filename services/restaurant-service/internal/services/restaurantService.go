package services

import (
	"context"
	"log"

	"github.com/suryansh74/zomato/services/restaurant-service/internal/models"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/repositories"
)

type RestaurantService struct {
	repo repositories.RestaurantRepository
}

func NewRestaurantService(repo repositories.RestaurantRepository) *RestaurantService {
	log.Println("Initializing RestaurantService")
	return &RestaurantService{repo: repo}
}

// ✅ Changed email to ownerID
func (s *RestaurantService) CheckIfOwnerHasRestaurant(ctx context.Context, ownerID string) (string, bool, error) {
	log.Println("Service: CheckIfOwnerHasRestaurant called with ownerID:", ownerID)

	id, exists, err := s.repo.CheckIfOwnerHasRestaurant(ctx, ownerID)

	log.Println("Service: CheckIfOwnerHasRestaurant result -> id:", id, "exists:", exists, "error:", err)
	return id, exists, err
}

// ✅ Changed email to ownerID and assigned to OwnerID field
func (s *RestaurantService) CreateRestaurant(ctx context.Context, ownerID string, req *models.RestaurantRequest) (*models.Restaurant, error) {
	log.Println("Service: CreateRestaurant called for owner:", ownerID, "name:", req.Name)

	restaurant := &models.Restaurant{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		OwnerID:     ownerID, // ✅ CHANGED HERE
		Phone:       req.Phone,
		IsVerified:  false,
		IsOpen:      false,
		AutoLocation: models.GeoJSONPoint{
			Type:             "Point",
			Coordinates:      []float64{req.Longitude, req.Latitude},
			FormattedAddress: req.FormattedAddress,
		},
	}

	result, err := s.repo.CreateRestaurant(ctx, restaurant)
	if err != nil {
		log.Println("Service: CreateRestaurant error:", err)
		return nil, err
	}

	log.Println("Service: Restaurant created with ID:", result.ID.Hex())
	return result, nil
}

// ✅ Changed email to ownerID
func (s *RestaurantService) GetRestaurant(ctx context.Context, ownerID string) (*models.Restaurant, error) {
	log.Println("Service: GetRestaurant called with ownerID:", ownerID)

	result, err := s.repo.GetRestaurant(ctx, ownerID)
	if err != nil {
		log.Println("Service: GetRestaurant error:", err)
		return nil, err
	}

	log.Println("Service: GetRestaurant success ID:", result.ID.Hex())
	return result, nil
}

// ✅ Changed email to ownerID
func (s *RestaurantService) UpdateRestaurant(ctx context.Context, ownerID string, req *models.UpdateRestaurantRequest) (*models.Restaurant, error) {
	log.Println("Service: UpdateRestaurant called for ownerID:", ownerID)

	result, err := s.repo.UpdateRestaurant(ctx, ownerID, req)
	if err != nil {
		log.Println("Service: UpdateRestaurant error:", err)
		return nil, err
	}

	log.Println("Service: UpdateRestaurant success ID:", result.ID.Hex())
	return result, nil
}

func (s *RestaurantService) GetNearbyRestaurants(ctx context.Context, lat, lon, radius float64, search string, isOpenFilter *bool) ([]models.Restaurant, error) {
	log.Println("Service: GetNearbyRestaurants called with lat:", lat, "lon:", lon, "radius:", radius, "search:", search, "isOpenFilter:", isOpenFilter)

	result, err := s.repo.GetNearbyRestaurants(ctx, lat, lon, radius, search, isOpenFilter)
	if err != nil {
		log.Println("Service: GetNearbyRestaurants error:", err)
		return nil, err
	}

	log.Println("Service: GetNearbyRestaurants count:", len(result))
	return result, nil
}

func (s *RestaurantService) GetRestaurantByID(ctx context.Context, id string) (*models.Restaurant, error) {
	log.Println("Service: GetRestaurantByID called with id:", id)

	result, err := s.repo.GetRestaurantByID(ctx, id)
	if err != nil {
		log.Println("Service: GetRestaurantByID error:", err)
		return nil, err
	}

	log.Println("Service: GetRestaurantByID success ID:", result.ID.Hex())
	return result, nil
}
