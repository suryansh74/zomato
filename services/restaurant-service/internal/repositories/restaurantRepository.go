package repositories

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RestaurantRepository interface{}

type restaurantRepository struct {
	db             *mongo.Client
	dbName         string
	collectionName string
}

func NewRestaurantRepository(db *mongo.Client, dbName, collectionName string) RestaurantRepository {
	return &restaurantRepository{
		db:             db,
		dbName:         dbName,
		collectionName: collectionName,
	}
}
