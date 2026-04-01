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
