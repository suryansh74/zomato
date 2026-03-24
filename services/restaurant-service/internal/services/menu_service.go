package services

import (
	"context"

	"github.com/suryansh74/zomato/services/restaurant-service/internal/models"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/repositories"
)

type MenuService struct {
	repo repositories.MenuRepository
}

func NewMenuService(repo repositories.MenuRepository) *MenuService {
	return &MenuService{repo: repo}
}

func (s *MenuService) CreateMenuItem(ctx context.Context, restaurantID string, req *models.MenuItemRequest) (*models.MenuItem, error) {
	item := &models.MenuItem{
		RestaurantID: restaurantID,
		Name:         req.Name,
		Description:  req.Description,
		Image:        req.Image,
		Price:        req.Price,
		IsAvailable:  req.IsAvailable,
	}
	return s.repo.CreateMenuItem(ctx, item)
}

func (s *MenuService) GetMenuItemsByRestaurant(ctx context.Context, restaurantID string) ([]models.MenuItem, error) {
	return s.repo.GetMenuItemsByRestaurant(ctx, restaurantID)
}

func (s *MenuService) GetMenuItemByID(ctx context.Context, id string, restaurantID string) (*models.MenuItem, error) {
	return s.repo.GetMenuItemByID(ctx, id, restaurantID)
}

func (s *MenuService) UpdateMenuItem(ctx context.Context, id string, restaurantID string, req *models.UpdateMenuItemRequest) (*models.MenuItem, error) {
	return s.repo.UpdateMenuItem(ctx, id, restaurantID, req)
}

func (s *MenuService) DeleteMenuItem(ctx context.Context, id string, restaurantID string) error {
	return s.repo.DeleteMenuItem(ctx, id, restaurantID)
}
