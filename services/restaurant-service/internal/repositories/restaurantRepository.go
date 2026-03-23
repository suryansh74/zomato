package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/suryansh74/zomato/services/restaurant-service/apperr"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RestaurantRepository interface {
	CheckIfOwnerHasRestaurant(ctx context.Context, email string) (bool, error)
	CreateRestaurant(ctx context.Context, restaurant *models.Restaurant) (*models.Restaurant, error) // ← add this
}

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

func (r *restaurantRepository) CheckIfOwnerHasRestaurant(ctx context.Context, email string) (bool, error) {
	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	var result models.Restaurant
	err := coll.FindOne(ctx, bson.D{{Key: "owner_email", Value: email}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return true, apperr.ErrInternalServer
	}
	return true, nil
}

func (r *restaurantRepository) CreateRestaurant(ctx context.Context, restaurant *models.Restaurant) (*models.Restaurant, error) {
	restaurant.CreatedAt = time.Now()
	restaurant.UpdatedAt = time.Now()

	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	result, err := coll.InsertOne(ctx, restaurant)
	if err != nil {
		return nil, apperr.ErrInternalServer
	}

	restaurant.ID = result.InsertedID.(bson.ObjectID)
	return restaurant, nil
}
