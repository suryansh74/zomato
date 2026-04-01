package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// OrderItem is a snapshot of the menu item at checkout
type OrderItem struct {
	ItemID   string  `json:"item_id"   bson:"item_id"`
	Name     string  `json:"name"      bson:"name"`
	Price    float64 `json:"price"     bson:"price"`
	Quantity int     `json:"quantity"  bson:"quantity"`
}

type Order struct {
	ID            bson.ObjectID `json:"id"              bson:"_id,omitempty"`
	UserID        string        `json:"user_id"         bson:"user_id"`
	RestaurantID  string        `json:"restaurant_id"   bson:"restaurant_id"`
	AddressID     string        `json:"address_id"      bson:"address_id"`
	Items         []OrderItem   `json:"items"           bson:"items"`
	ItemTotal     float64       `json:"item_total"      bson:"item_total"`
	PlatformFee   float64       `json:"platform_fee"    bson:"platform_fee"`
	DeliveryFee   float64       `json:"delivery_fee"    bson:"delivery_fee"`
	GrandTotal    float64       `json:"grand_total"     bson:"grand_total"`
	Status        string        `json:"status"          bson:"status"`         // unpaid, paid, preparing, delivered
	PaymentMethod string        `json:"payment_method"  bson:"payment_method"` // stripe
	CreatedAt     time.Time     `json:"created_at"      bson:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"      bson:"updated_at"`
}

type CreateOrderRequest struct {
	AddressID string `json:"address_id" validate:"required"`
}
