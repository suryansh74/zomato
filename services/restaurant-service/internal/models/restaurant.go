package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Restaurant struct {
	ID           bson.ObjectID `json:"id"                bson:"_id,omitempty"`
	Name         string        `json:"name"              bson:"name"             validate:"required,min=2,max=100"`
	Description  string        `json:"description"       bson:"description"      validate:"omitempty,max=500"`
	Image        string        `json:"image"             bson:"image"            validate:"required,url"`
	OwnerID      string        `json:"owner_id"          bson:"owner_id"         validate:"required"`
	Phone        string        `json:"phone"             bson:"phone"            validate:"required,e164"`
	IsVerified   bool          `json:"is_verified"       bson:"is_verified"`
	AutoLocation GeoJSONPoint  `json:"auto_location"     bson:"auto_location"    validate:"required"`
	IsOpen       bool          `json:"is_open"           bson:"is_open"`
	DistanceKm   *float64      `json:"distance_km,omitempty" bson:"distanceKm,omitempty"`
	CreatedAt    time.Time     `json:"created_at"        bson:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"        bson:"updated_at"`
}

type GeoJSONPoint struct {
	Type             string    `json:"type"              bson:"type"              validate:"required,oneof=Point"`
	Coordinates      []float64 `json:"coordinates"       bson:"coordinates"       validate:"required,len=2,dive,number"`
	FormattedAddress string    `json:"formatted_address" bson:"formatted_address" validate:"required,min=5,max=200"`
}

type RestaurantRequest struct {
	Name             string  `json:"name"        bson:"name"         validate:"required,min=2,max=100"`
	Description      string  `json:"description" bson:"description"  validate:"omitempty,max=500"`
	Image            string  `json:"image"       bson:"image"        validate:"required,url"`
	Phone            string  `json:"phone"       bson:"phone"        validate:"required,e164"`
	Latitude         float64 `json:"latitude"    bson:"latitude"     validate:"required"`
	Longitude        float64 `json:"longitude"   bson:"longitude"    validate:"required"`
	FormattedAddress string  `json:"formatted_address" bson:"formatted_address" validate:"required,min=5,max=200"`
}

type UpdateRestaurantRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	IsOpen      *bool   `json:"is_open"`
}

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
