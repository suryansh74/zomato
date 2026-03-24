package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/suryansh74/zomato/services/restaurant-service/apperr"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RestaurantRepository interface {
	CheckIfOwnerHasRestaurant(ctx context.Context, email string) (string, bool, error)
	CreateRestaurant(ctx context.Context, restaurant *models.Restaurant) (*models.Restaurant, error)
	GetRestaurant(ctx context.Context, email string) (*models.Restaurant, error)
	UpdateRestaurant(ctx context.Context, email string, req *models.UpdateRestaurantRequest) (*models.Restaurant, error)
}

type restaurantRepository struct {
	db             *mongo.Client
	dbName         string
	collectionName string
}

// NewRestaurantService constructor for restaurantService
// ================================================================================
func NewRestaurantRepository(db *mongo.Client, dbName, collectionName string) RestaurantRepository {
	return &restaurantRepository{
		db:             db,
		dbName:         dbName,
		collectionName: collectionName,
	}
}

// NewRestaurantService constructor for restaurantService
// ================================================================================
func (r *restaurantRepository) CheckIfOwnerHasRestaurant(ctx context.Context, email string) (string, bool, error) {
	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	var result models.Restaurant
	err := coll.FindOne(ctx, bson.D{{Key: "owner_email", Value: email}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", false, nil
		}
		return "", true, apperr.ErrInternalServer
	}
	return result.ID.String(), true, nil
}

// NewRestaurantService constructor for restaurantService
// ================================================================================
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

// GetRestaurant
// ================================================================================
func (r *restaurantRepository) GetRestaurant(ctx context.Context, email string) (*models.Restaurant, error) {
	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	var result models.Restaurant
	err := coll.FindOne(ctx, bson.D{{Key: "owner_email", Value: email}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperr.ErrRestaurantNotFound
		}
		return nil, apperr.ErrInternalServer
	}
	return &result, nil
}

func (r *restaurantRepository) UpdateRestaurant(ctx context.Context, email string, req *models.UpdateRestaurantRequest) (*models.Restaurant, error) {
	coll := r.db.Database(r.dbName).Collection(r.collectionName)

	updateDocs := bson.M{}
	if req.Name != nil {
		updateDocs["name"] = *req.Name
	}
	if req.Description != nil {
		updateDocs["description"] = *req.Description
	}
	if req.IsOpen != nil {
		updateDocs["is_open"] = *req.IsOpen
	}
	updateDocs["updated_at"] = time.Now()

	update := bson.M{"$set": updateDocs}

	// Return the updated document
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated models.Restaurant
	err := coll.FindOneAndUpdate(ctx, bson.D{{Key: "owner_email", Value: email}}, update, opts).Decode(&updated)
	if err != nil {
		return nil, apperr.ErrInternalServer
	}

	return &updated, nil
}
