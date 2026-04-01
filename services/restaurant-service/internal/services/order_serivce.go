package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/models"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/repositories"
)

type OrderService struct {
	repo            repositories.OrderRepository
	cartSrv         *CartService
	menuSrv         *MenuService
	stripeSecretKey string
}

func NewOrderService(repo repositories.OrderRepository, cartSrv *CartService, menuSrv *MenuService, stripeSecretKey string) *OrderService {
	log.Println("Initializing OrderService")
	return &OrderService{
		repo:            repo,
		cartSrv:         cartSrv,
		menuSrv:         menuSrv,
		stripeSecretKey: stripeSecretKey,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, req *models.CreateOrderRequest) (*models.Order, error) {
	cartItems, err := s.cartSrv.GetCartByUserID(ctx, userID)
	if err != nil || len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	var orderItems []models.OrderItem
	var itemTotal float64 = 0
	restaurantID := cartItems[0].RestaurantID

	for _, cartItem := range cartItems {
		menuItem, err := s.menuSrv.GetMenuItemByID(ctx, cartItem.ItemID, restaurantID)
		if err != nil {
			continue
		}

		itemTotal += menuItem.Price * float64(cartItem.Quantity)
		orderItems = append(orderItems, models.OrderItem{
			ItemID:   menuItem.ID.Hex(),
			Name:     menuItem.Name,
			Price:    menuItem.Price,
			Quantity: cartItem.Quantity,
		})
	}

	if len(orderItems) == 0 {
		return nil, errors.New("no valid items in cart")
	}

	platformFee := 5.0
	deliveryFee := 0.0
	grandTotal := itemTotal + platformFee + deliveryFee

	order := &models.Order{
		UserID:        userID,
		RestaurantID:  restaurantID,
		AddressID:     req.AddressID,
		Items:         orderItems,
		ItemTotal:     itemTotal,
		PlatformFee:   platformFee,
		DeliveryFee:   deliveryFee,
		GrandTotal:    grandTotal,
		Status:        "unpaid",
		PaymentMethod: "stripe",
	}

	createdOrder, err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	// Clear the user's cart after successfully drafting the order
	_ = s.cartSrv.ClearCart(ctx, userID)

	return createdOrder, nil
}

func (s *OrderService) CreateStripeSession(ctx context.Context, orderID string, userID string) (string, error) {
	stripe.Key = s.stripeSecretKey

	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return "", err
	}

	if order.UserID != userID {
		return "", errors.New("unauthorized access to order")
	}
	if order.Status == "paid" {
		return "", errors.New("order is already paid")
	}

	var lineItems []*stripe.CheckoutSessionLineItemParams
	for _, item := range order.Items {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("inr"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(item.Name),
				},
				UnitAmount: stripe.Int64(int64(item.Price * 100)), // Stripe expects paise for INR
			},
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	// Add platform fee
	lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
		PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
			Currency: stripe.String("inr"),
			ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
				Name: stripe.String("Platform Fee"),
			},
			UnitAmount: stripe.Int64(int64(order.PlatformFee * 100)),
		},
		Quantity: stripe.Int64(1),
	})

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String("http://localhost:5173/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:          stripe.String("http://localhost:5173/cart"),
		Metadata: map[string]string{
			"order_id": orderID, // Crucial for the webhook to identify the order
		},
	}

	sess, err := session.New(params)
	if err != nil {
		log.Println("Stripe session creation failed:", err)
		return "", err
	}

	return sess.URL, nil
}

func (s *OrderService) ProcessPaymentSuccess(ctx context.Context, orderID string) (*models.Order, error) {
	log.Println("Processing successful payment for Order:", orderID)
	return s.repo.MarkOrderAsPaid(ctx, orderID)
}

func (s *OrderService) NotifyRestaurant(order *models.Order) {
	// 1. Prepare the payload we want to send to the dashboard
	payload := map[string]interface{}{
		"restaurant_id": order.RestaurantID,
		"order_id":      order.ID.Hex(),
		"items":         order.Items,
		"grand_total":   order.GrandTotal,
		"status":        order.Status,
		"message":       "New Order Received!",
	}

	jsonData, _ := json.Marshal(payload)

	// 2. Create the HTTP request to your realtime-service
	req, err := http.NewRequest("POST", "http://localhost:8003/api/internal/notify-order", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Failed to create internal request:", err)
		return
	}

	// 3. Attach the VIP Internal Key
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Internal-Key", "super_secret_zomato_internal_key_2026")

	// 4. Fire the request in the background
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to notify realtime service:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Realtime service returned status: %d", resp.StatusCode)
	} else {
		log.Println("✅ Successfully notified Realtime Service about Order:", order.ID.Hex())
	}
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID string, restaurantID string, newStatus string) error {
	// First, verify this order actually belongs to this restaurant!
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.RestaurantID != restaurantID {
		return errors.New("unauthorized: order does not belong to your restaurant")
	}

	return s.repo.UpdateOrderStatus(ctx, orderID, newStatus)
}

// Add to bottom of order_service.go
func (s *OrderService) GetActiveOrders(ctx context.Context, restaurantID string) ([]models.Order, error) {
	return s.repo.GetActiveOrdersByRestaurant(ctx, restaurantID)
}
