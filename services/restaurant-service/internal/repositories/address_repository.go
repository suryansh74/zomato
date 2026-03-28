package repositories

import (
	"context"
	"log"
	"time"

	"github.com/suryansh74/zomato/services/restaurant-service/apperr"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AddressRepository interface {
	CreateAddress(ctx context.Context, address *models.Address) (*models.Address, error)
	GetAddressesByUserID(ctx context.Context, userID string) ([]models.Address, error)
	DeleteAddressByID(ctx context.Context, id string, userID string) error
}

type addressRepository struct {
	db             *mongo.Client
	dbName         string
	collectionName string
}

func NewAddressRepository(db *mongo.Client, dbName, collectionName string) AddressRepository {
	log.Println("Initializing AddressRepository")
	return &addressRepository{
		db:             db,
		dbName:         dbName,
		collectionName: collectionName,
	}
}

func (r *addressRepository) CreateAddress(ctx context.Context, address *models.Address) (*models.Address, error) {
	address.CreatedAt = time.Now()
	address.UpdatedAt = time.Now()

	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	result, err := coll.InsertOne(ctx, address)
	if err != nil {
		log.Println("Error inserting address:", err)
		return nil, apperr.ErrInternalServer
	}

	address.ID = result.InsertedID.(bson.ObjectID)
	return address, nil
}

func (r *addressRepository) GetAddressesByUserID(ctx context.Context, userID string) ([]models.Address, error) {
	coll := r.db.Database(r.dbName).Collection(r.collectionName)

	cursor, err := coll.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		log.Println("Error fetching addresses:", err)
		return nil, apperr.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var addresses []models.Address
	if err = cursor.All(ctx, &addresses); err != nil {
		log.Println("Error decoding addresses:", err)
		return nil, apperr.ErrInternalServer
	}

	if addresses == nil {
		addresses = []models.Address{}
	}

	return addresses, nil
}

func (r *addressRepository) DeleteAddressByID(ctx context.Context, id string, userID string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return apperr.ErrInvalidID
	}

	coll := r.db.Database(r.dbName).Collection(r.collectionName)

	// Ensure the user deleting the address actually owns it!
	result, err := coll.DeleteOne(ctx, bson.M{"_id": objID, "user_id": userID})
	if err != nil {
		log.Println("Error deleting address:", err)
		return apperr.ErrInternalServer
	}

	if result.DeletedCount == 0 {
		return apperr.ErrMenuNotFound // You can reuse this or create ErrAddressNotFound
	}

	return nil
}
