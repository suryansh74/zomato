package services

import (
	"context"
	"log"

	"github.com/suryansh74/zomato/services/restaurant-service/internal/models"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/repositories"
)

type MenuService struct {
	repo repositories.MenuRepository
}

func NewMenuService(repo repositories.MenuRepository) *MenuService {
	log.Println("Initializing MenuService")
	return &MenuService{repo: repo}
}

func (s *MenuService) CreateMenuItem(ctx context.Context, restaurantID string, req *models.MenuItemRequest) (*models.MenuItem, error) {
	log.Println("Service: CreateMenuItem called for restaurantID:", restaurantID, "name:", req.Name)

	item := &models.MenuItem{
		RestaurantID: restaurantID,
		Name:         req.Name,
		Description:  req.Description,
		Image:        req.Image,
		Price:        req.Price,
		IsAvailable:  req.IsAvailable,
	}

	result, err := s.repo.CreateMenuItem(ctx, item)
	if err != nil {
		log.Println("Service: CreateMenuItem error:", err)
		return nil, err
	}

	log.Println("Service: Menu item created with ID:", result.ID.Hex())
	return result, nil
}

func (s *MenuService) GetMenuItemsByRestaurant(ctx context.Context, restaurantID string) ([]models.MenuItem, error) {
	log.Println("Service: GetMenuItemsByRestaurant called for restaurantID:", restaurantID)

	result, err := s.repo.GetMenuItemsByRestaurant(ctx, restaurantID)
	if err != nil {
		log.Println("Service: GetMenuItemsByRestaurant error:", err)
		return nil, err
	}

	log.Println("Service: Menu items fetched count:", len(result))
	return result, nil
}

func (s *MenuService) GetMenuItemByID(ctx context.Context, id string, restaurantID string) (*models.MenuItem, error) {
	log.Println("Service: GetMenuItemByID called with id:", id, "restaurantID:", restaurantID)

	result, err := s.repo.GetMenuItemByID(ctx, id, restaurantID)
	if err != nil {
		log.Println("Service: GetMenuItemByID error:", err)
		return nil, err
	}

	log.Println("Service: Menu item fetched ID:", result.ID.Hex())
	return result, nil
}

func (s *MenuService) UpdateMenuItem(ctx context.Context, id string, restaurantID string, req *models.UpdateMenuItemRequest) (*models.MenuItem, error) {
	log.Println("Service: UpdateMenuItem called with id:", id, "restaurantID:", restaurantID)

	result, err := s.repo.UpdateMenuItem(ctx, id, restaurantID, req)
	if err != nil {
		log.Println("Service: UpdateMenuItem error:", err)
		return nil, err
	}

	if result != nil {
		log.Println("Service: Menu item updated ID:", result.ID.Hex())
	} else {
		log.Println("Service: No update performed for ID:", id)
	}

	return result, nil
}

func (s *MenuService) DeleteMenuItem(ctx context.Context, id string, restaurantID string) error {
	log.Println("Service: DeleteMenuItem called with id:", id, "restaurantID:", restaurantID)

	err := s.repo.DeleteMenuItem(ctx, id, restaurantID)
	if err != nil {
		log.Println("Service: DeleteMenuItem error:", err)
		return err
	}

	log.Println("Service: Menu item deleted ID:", id)
	return nil
}
