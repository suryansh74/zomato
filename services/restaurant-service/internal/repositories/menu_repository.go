package repositories

import (
	"context"
	"log"
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
	log.Println("Initializing MenuRepository")
	return &menuRepository{
		db:             db,
		dbName:         dbName,
		collectionName: collectionName,
	}
}

func (r *menuRepository) CreateMenuItem(ctx context.Context, item *models.MenuItem) (*models.MenuItem, error) {
	log.Println("Repo: CreateMenuItem called for:", item.Name)

	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	result, err := coll.InsertOne(ctx, item)
	if err != nil {
		log.Println("Error inserting menu item:", err)
		return nil, apperr.ErrInternalServer
	}

	item.ID = result.InsertedID.(bson.ObjectID)
	log.Println("Menu item created with ID:", item.ID.Hex())

	return item, nil
}

func (r *menuRepository) GetMenuItemsByRestaurant(ctx context.Context, restaurantID string) ([]models.MenuItem, error) {
	log.Println("Repo: GetMenuItemsByRestaurant called for restaurantID:", restaurantID)

	coll := r.db.Database(r.dbName).Collection(r.collectionName)

	cursor, err := coll.Find(ctx, bson.D{{Key: "restaurant_id", Value: restaurantID}})
	if err != nil {
		log.Println("Error fetching menu items:", err)
		return nil, apperr.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var items []models.MenuItem
	if err = cursor.All(ctx, &items); err != nil {
		log.Println("Error decoding menu items:", err)
		return nil, apperr.ErrInternalServer
	}

	if items == nil {
		items = []models.MenuItem{}
	}

	log.Println("Menu items fetched count:", len(items))
	return items, nil
}

func (r *menuRepository) GetMenuItemByID(ctx context.Context, id string, restaurantID string) (*models.MenuItem, error) {
	log.Println("Repo: GetMenuItemByID called with id:", id, "restaurantID:", restaurantID)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid ObjectID:", err)
		return nil, apperr.ErrInvalidID
	}

	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	var item models.MenuItem

	err = coll.FindOne(ctx, bson.M{"_id": objID, "restaurant_id": restaurantID}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Menu item not found for id:", id)
			return nil, apperr.ErrMenuNotFound
		}
		log.Println("Error fetching menu item:", err)
		return nil, apperr.ErrInternalServer
	}

	log.Println("Menu item fetched:", item.ID.Hex())
	return &item, nil
}

func (r *menuRepository) UpdateMenuItem(ctx context.Context, id string, restaurantID string, req *models.UpdateMenuItemRequest) (*models.MenuItem, error) {
	log.Println("Repo: UpdateMenuItem called with id:", id, "restaurantID:", restaurantID)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid ObjectID:", err)
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
		log.Println("No fields to update for id:", id)
		return nil, nil
	}

	updateDocs["updated_at"] = time.Now()
	log.Println("Update fields:", updateDocs)

	update := bson.M{"$set": updateDocs}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.MenuItem

	err = coll.FindOneAndUpdate(ctx, bson.M{"_id": objID, "restaurant_id": restaurantID}, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Menu item not found for update:", id)
			return nil, apperr.ErrMenuNotFound
		}
		log.Println("Error updating menu item:", err)
		return nil, apperr.ErrInternalServer
	}

	log.Println("Menu item updated:", updated.ID.Hex())
	return &updated, nil
}

func (r *menuRepository) DeleteMenuItem(ctx context.Context, id string, restaurantID string) error {
	log.Println("Repo: DeleteMenuItem called with id:", id, "restaurantID:", restaurantID)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid ObjectID:", err)
		return apperr.ErrInvalidID
	}

	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	result, err := coll.DeleteOne(ctx, bson.M{"_id": objID, "restaurant_id": restaurantID})
	if err != nil {
		log.Println("Error deleting menu item:", err)
		return apperr.ErrInternalServer
	}

	if result.DeletedCount == 0 {
		log.Println("Menu item not found for deletion:", id)
		return apperr.ErrMenuNotFound
	}

	log.Println("Menu item deleted successfully:", id)
	return nil
}
