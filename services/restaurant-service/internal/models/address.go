package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Address represents the database schema matching the instructor's TS interface
type Address struct {
	ID               bson.ObjectID `json:"id"                bson:"_id,omitempty"`
	UserID           string        `json:"user_id"           bson:"user_id"`
	Mobile           string        `json:"mobile"            bson:"mobile"` // Using string for safety with phone numbers
	FormattedAddress string        `json:"formatted_address" bson:"formatted_address"`
	Location         GeoJSONPoint  `json:"location"          bson:"location"`
	CreatedAt        time.Time     `json:"created_at"        bson:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"        bson:"updated_at"`
}

// AddressRequest is what the frontend will send us
type AddressRequest struct {
	Mobile           string  `json:"mobile"            validate:"required"`
	FormattedAddress string  `json:"formatted_address" validate:"required"`
	Longitude        float64 `json:"longitude"         validate:"required"`
	Latitude         float64 `json:"latitude"          validate:"required"`
}
