package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `json:"id"         bson:"_id,omitempty"`
	Name      string        `json:"name"       bson:"name"          validate:"required,min=2,max=50"`
	Email     string        `json:"email"      bson:"email"         validate:"required,email"`
	Image     string        `json:"image"      bson:"image"         validate:"required"`
	Role      string        `json:"role"       bson:"role"          validate:"omitempty,oneof=customer restaurant_owner admin"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

type LoginRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=50"`
	Email string `json:"email"    validate:"required,email"`
	Image string `json:"image"    validate:"required"`
}
