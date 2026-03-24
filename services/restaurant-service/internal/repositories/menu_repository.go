package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/suryansh74/zomato/services/restaurant-service/apperr"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/models"
)

type MenuRepository interface {
	CreateMenuItem(ctx context.Context, item *models.MenuItem) (*models.MenuItem, error)
	GetMenuItemsByRestaurant(ctx context.Context, restaurantID string) ([]models.MenuItem, error)
	GetMenuItemByID(ctx context.Context, id string, restaurantID string) (*models.MenuItem, error)
	UpdateMenuItem(ctx context.Context, id string, restaurantID string, req *models.UpdateMenuItemRequest) (*models.MenuItem, error)
	DeleteMenuItem(ctx context.Context, id string, restaurantID string) error
}

type menuRepository struct {
	db             *mongo.Client
	dbName         string
	collectionName string
}

func NewMenuRepository(db *mongo.Client, dbName, collectionName string) MenuRepository {
	return &menuRepository{
		db:             db,
		dbName:         dbName,
		collectionName: collectionName,
	}
}

func (r *menuRepository) CreateMenuItem(ctx context.Context, item *models.MenuItem) (*models.MenuItem, error) {
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	result, err := coll.InsertOne(ctx, item)
	if err != nil {
		return nil, apperr.ErrInternalServer
	}

	item.ID = result.InsertedID.(bson.ObjectID)
	return item, nil
}

func (r *menuRepository) GetMenuItemsByRestaurant(ctx context.Context, restaurantID string) ([]models.MenuItem, error) {
	coll := r.db.Database(r.dbName).Collection(r.collectionName)

	cursor, err := coll.Find(ctx, bson.D{{Key: "restaurant_id", Value: restaurantID}})
	if err != nil {
		return nil, apperr.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var items []models.MenuItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, apperr.ErrInternalServer
	}

	// Return empty slice instead of null if no items found
	if items == nil {
		items = []models.MenuItem{}
	}

	return items, nil
}

func (r *menuRepository) GetMenuItemByID(ctx context.Context, id string, restaurantID string) (*models.MenuItem, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperr.ErrInvalidID
	}

	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	var item models.MenuItem
	err = coll.FindOne(ctx, bson.M{"_id": objID, "restaurant_id": restaurantID}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperr.ErrMenuNotFound
		}
		return nil, apperr.ErrInternalServer
	}
	return &item, nil
}

func (r *menuRepository) UpdateMenuItem(ctx context.Context, id string, restaurantID string, req *models.UpdateMenuItemRequest) (*models.MenuItem, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperr.ErrInvalidID
	}

	coll := r.db.Database(r.dbName).Collection(r.collectionName)

	updateDocs := bson.M{}
	if req.Name != nil {
		updateDocs["name"] = *req.Name
	}
	if req.Description != nil {
		updateDocs["description"] = *req.Description
	}
	if req.Image != nil {
		updateDocs["image"] = *req.Image
	}
	if req.Price != nil {
		updateDocs["price"] = *req.Price
	}
	if req.IsAvailable != nil {
		updateDocs["is_available"] = *req.IsAvailable
	}

	if len(updateDocs) == 0 {
		return nil, nil // Nothing to update
	}

	updateDocs["updated_at"] = time.Now()
	update := bson.M{"$set": updateDocs}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.MenuItem

	err = coll.FindOneAndUpdate(ctx, bson.M{"_id": objID, "restaurant_id": restaurantID}, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperr.ErrMenuNotFound
		}
		return nil, apperr.ErrInternalServer
	}

	return &updated, nil
}

func (r *menuRepository) DeleteMenuItem(ctx context.Context, id string, restaurantID string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return apperr.ErrInvalidID
	}

	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	result, err := coll.DeleteOne(ctx, bson.M{"_id": objID, "restaurant_id": restaurantID})
	if err != nil {
		return apperr.ErrInternalServer
	}

	if result.DeletedCount == 0 {
		return apperr.ErrMenuNotFound
	}

	return nil
}
