package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MenuItem struct {
	ID           bson.ObjectID `json:"id"             bson:"_id,omitempty"`
	RestaurantID string        `json:"restaurant_id"  bson:"restaurant_id"`
	Name         string        `json:"name"           bson:"name"          validate:"required,min=2,max=100"`
	Description  string        `json:"description"    bson:"description"   validate:"omitempty,max=500"`
	Image        string        `json:"image"          bson:"image"         validate:"omitempty,url"`
	Price        float64       `json:"price"          bson:"price"         validate:"required"`
	IsAvailable  bool          `json:"is_available"   bson:"is_available"  validate:"required"`
	CreatedAt    time.Time     `json:"created_at"     bson:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"     bson:"updated_at"`
}

type MenuItemRequest struct {
	Name        string  `json:"name"          validate:"required,min=2,max=100"`
	Description string  `json:"description"   validate:"omitempty,max=500"`
	Image       string  `json:"image"         validate:"omitempty,url"`
	Price       float64 `json:"price"         validate:"required"`
	IsAvailable bool    `json:"is_available"`
}

type UpdateMenuItemRequest struct {
	Name        *string  `json:"name"          validate:"omitempty,min=2,max=100"`
	Description *string  `json:"description"   validate:"omitempty,max=500"`
	Image       *string  `json:"image"         validate:"omitempty,url"`
	Price       *float64 `json:"price"`
	IsAvailable *bool    `json:"is_available"`
}
