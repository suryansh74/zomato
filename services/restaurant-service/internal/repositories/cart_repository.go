package repositories

import (
	"context"
	"log"
	"time"

	"github.com/suryansh74/zomato/services/restaurant-service/apperr"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CartRepository interface {
	AddToCart(ctx context.Context, req *models.CartRequest) (*models.Cart, error)
	GetCartByUserID(ctx context.Context, userID string) ([]models.Cart, error)
	UpdateQuantity(ctx context.Context, userID, itemID, action string) error
	ClearCart(ctx context.Context, userID string) error
}

type cartRepository struct {
	db             *mongo.Client
	dbName         string
	collectionName string
}

func NewCartRepository(db *mongo.Client, dbName, collectionName string) CartRepository {
	log.Println("Initializing MenuRepository")
	return &cartRepository{
		db:             db,
		dbName:         dbName,
		collectionName: collectionName,
	}
}

// internal/repositories/cartRepository.go
func (r *cartRepository) AddToCart(ctx context.Context, req *models.CartRequest) (*models.Cart, error) {
	coll := r.db.Database(r.dbName).Collection(r.collectionName)

	// 1. CONFLICT CHECK: Does this user have ANY items in their cart from a DIFFERENT restaurant?
	conflictFilter := bson.M{
		"user_id":       req.UserID,
		"restaurant_id": bson.M{"$ne": req.RestaurantID},
	}

	// We just need to know if ONE document exists that matches the conflict
	var existingItem models.Cart
	err := coll.FindOne(ctx, conflictFilter).Decode(&existingItem)
	if err != mongo.ErrNoDocuments {
		// If we found a document (err == nil) OR there was a database error, we have a conflict/issue.
		if err == nil {
			log.Println("Cart conflict: User is trying to order from multiple restaurants")
			return nil, apperr.ErrCartConflict
		}
		return nil, apperr.ErrInternalServer
	}

	// 2. UPSERT ITEM: If no conflict, we add the item (or increment its quantity)
	filter := bson.M{
		"user_id":       req.UserID,
		"restaurant_id": req.RestaurantID,
		"item_id":       req.ItemID,
	}

	update := bson.M{
		"$inc": bson.M{"quantity": 1},
		// setOnInsert only runs if the document is brand new.
		// Time needs to be explicitly set here if your struct doesn't auto-handle it on DB insertion
		"$setOnInsert": bson.M{
			"user_id":       req.UserID,
			"restaurant_id": req.RestaurantID,
			"item_id":       req.ItemID,
			"created_at":    time.Now(),
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var updatedCart models.Cart
	err = coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedCart)
	if err != nil {
		log.Println("Error upserting cart item:", err)
		return nil, apperr.ErrInternalServer
	}

	return &updatedCart, nil
}

func (r *cartRepository) GetCartByUserID(ctx context.Context, userID string) ([]models.Cart, error) {
	log.Println("Repo: GetCartByUserID called for user:", userID)

	coll := r.db.Database(r.dbName).Collection(r.collectionName)

	// ✅ FIX: Change "userId" to "user_id" to match MongoDB exactly
	cursor, err := coll.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		log.Println("Error fetching cart items:", err)
		return nil, apperr.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var cartItems []models.Cart
	if err = cursor.All(ctx, &cartItems); err != nil {
		log.Println("Error decoding cart items:", err)
		return nil, apperr.ErrInternalServer
	}

	if cartItems == nil {
		cartItems = []models.Cart{}
	}

	return cartItems, nil
}

func (r *cartRepository) UpdateQuantity(ctx context.Context, userID, itemID, action string) error {
	coll := r.db.Database(r.dbName).Collection(r.collectionName)

	// 1. Find the current item
	var item models.Cart
	err := coll.FindOne(ctx, bson.M{"user_id": userID, "item_id": itemID}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return apperr.ErrMenuNotFound
		}
		return apperr.ErrInternalServer
	}

	// 2. If decrementing and quantity is 1, delete it completely
	if action == "dec" && item.Quantity <= 1 {
		_, err = coll.DeleteOne(ctx, bson.M{"_id": item.ID})
		return err
	}

	// 3. Otherwise, increment or decrement
	incVal := 1
	if action == "dec" {
		incVal = -1
	}

	_, err = coll.UpdateOne(
		ctx,
		bson.M{"_id": item.ID},
		bson.M{
			"$inc": bson.M{"quantity": incVal},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)
	return err
}

func (r *cartRepository) ClearCart(ctx context.Context, userID string) error {
	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	_, err := coll.DeleteMany(ctx, bson.M{"user_id": userID})
	return err
}
