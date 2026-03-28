package services

import (
	"context"
	"log"

	"github.com/suryansh74/zomato/services/restaurant-service/internal/models"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/repositories"
)

type AddressService struct {
	repo repositories.AddressRepository
}

func NewAddressService(repo repositories.AddressRepository) *AddressService {
	log.Println("Initializing AddressService")
	return &AddressService{repo: repo}
}

func (s *AddressService) CreateAddress(ctx context.Context, userID string, req *models.AddressRequest) (*models.Address, error) {
	address := &models.Address{
		UserID:           userID,
		Mobile:           req.Mobile,
		FormattedAddress: req.FormattedAddress,
		Location: models.GeoJSONPoint{
			Type:        "Point",
			Coordinates: []float64{req.Longitude, req.Latitude}, // Remember: Longitude always comes first!
		},
	}

	return s.repo.CreateAddress(ctx, address)
}

func (s *AddressService) GetAddressesByUserID(ctx context.Context, userID string) ([]models.Address, error) {
	return s.repo.GetAddressesByUserID(ctx, userID)
}

func (s *AddressService) DeleteAddress(ctx context.Context, id string, userID string) error {
	return s.repo.DeleteAddressByID(ctx, id, userID)
}
