package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CartRequest struct {
	UserID       string `json:"-"` // We fill this from the token, not the JSON body
	RestaurantID string `json:"restaurantId" validate:"required"`
	ItemID       string `json:"itemId" validate:"required"`
}

type Cart struct {
	ID           bson.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID       string        `json:"userId" bson:"user_id"`
	RestaurantID string        `json:"restaurantId" bson:"restaurant_id"`
	ItemID       string        `json:"itemId" bson:"item_id"`
	Quantity     int           `json:"quantity" bson:"quantity"`
	CreatedAt    time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" bson:"updated_at"`
}

type CartUpdateRequest struct {
	ItemID string `json:"itemId" validate:"required"`
	Action string `json:"action" validate:"required,oneof=inc dec"` // Must be "inc" or "dec"
}
