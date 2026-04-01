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

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error)
	GetOrderByID(ctx context.Context, orderID string) (*models.Order, error)
	MarkOrderAsPaid(ctx context.Context, orderID string) (*models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status string) error

	GetActiveOrdersByRestaurant(ctx context.Context, restaurantID string) ([]models.Order, error)
}

type orderRepository struct {
	db             *mongo.Client
	dbName         string
	collectionName string
}

func NewOrderRepository(db *mongo.Client, dbName, collectionName string) OrderRepository {
	log.Println("Initializing OrderRepository")
	return &orderRepository{
		db:             db,
		dbName:         dbName,
		collectionName: collectionName,
	}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	result, err := coll.InsertOne(ctx, order)
	if err != nil {
		log.Println("Error inserting order:", err)
		return nil, apperr.ErrInternalServer
	}

	order.ID = result.InsertedID.(bson.ObjectID)
	return order, nil
}

func (r *orderRepository) GetOrderByID(ctx context.Context, orderID string) (*models.Order, error) {
	objID, err := bson.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, apperr.ErrInvalidID
	}

	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	var order models.Order

	err = coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperr.ErrMenuNotFound
		}
		return nil, apperr.ErrInternalServer
	}
	return &order, nil
}

func (r *orderRepository) MarkOrderAsPaid(ctx context.Context, orderID string) (*models.Order, error) {
	objID, err := bson.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, apperr.ErrInvalidID
	}

	coll := r.db.Database(r.dbName).Collection(r.collectionName)

	update := bson.M{
		"$set": bson.M{
			"status":     "paid",
			"updated_at": time.Now(),
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedOrder models.Order
	err = coll.FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&updatedOrder)
	if err != nil {
		log.Println("Error updating order to paid:", err)
		return nil, apperr.ErrInternalServer
	}

	log.Println("Order marked as paid:", orderID)
	return &updatedOrder, nil
}

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	objID, err := bson.ObjectIDFromHex(orderID)
	if err != nil {
		return apperr.ErrInvalidID
	}

	coll := r.db.Database(r.dbName).Collection(r.collectionName)

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := coll.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		log.Println("Error updating order status:", err)
		return apperr.ErrInternalServer
	}

	if result.MatchedCount == 0 {
		return apperr.ErrMenuNotFound // You can create a specific ErrOrderNotFound later
	}

	return nil
}

func (r *orderRepository) GetActiveOrdersByRestaurant(ctx context.Context, restaurantID string) ([]models.Order, error) {
	coll := r.db.Database(r.dbName).Collection(r.collectionName)

	// Filter: Find orders for this restaurant where status is "paid" or "preparing"
	filter := bson.M{
		"restaurant_id": restaurantID,
		"status":        bson.M{"$in": []string{"paid", "preparing"}},
	}

	// Sort by newest first
	opts := options.Find().SetSort(bson.M{"created_at": -1})

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
