package services

import (
	"context"
	"log"

	"github.com/suryansh74/zomato/services/restaurant-service/internal/models"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/repositories"
)

type CartService struct {
	repo repositories.CartRepository
}

func NewCartService(repo repositories.CartRepository) *CartService {
	log.Println("Initializing MenuService")
	return &CartService{repo: repo}
}

func (s *CartService) AddToCart(ctx context.Context, req *models.CartRequest) (*models.Cart, error) {
	log.Println("Service: AddToCart called")

	return s.repo.AddToCart(ctx, req)
}

func (s *CartService) GetCartByUserID(ctx context.Context, userID string) ([]models.Cart, error) {
	log.Println("Service: GetCartByUserID called")
	return s.repo.GetCartByUserID(ctx, userID)
}

func (s *CartService) UpdateQuantity(ctx context.Context, userID, itemID, action string) error {
	return s.repo.UpdateQuantity(ctx, userID, itemID, action)
}

func (s *CartService) ClearCart(ctx context.Context, userID string) error {
	return s.repo.ClearCart(ctx, userID)
}
