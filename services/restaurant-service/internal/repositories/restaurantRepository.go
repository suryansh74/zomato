package repositories

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/suryansh74/zomato/services/restaurant-service/apperr"
	"github.com/suryansh74/zomato/services/restaurant-service/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RestaurantRepository interface {
	CheckIfOwnerHasRestaurant(ctx context.Context, ownerID string) (string, bool, error) // ✅ Changed to ownerID
	CreateRestaurant(ctx context.Context, restaurant *models.Restaurant) (*models.Restaurant, error)
	GetRestaurant(ctx context.Context, ownerID string) (*models.Restaurant, error)                                         // ✅ Changed to ownerID
	UpdateRestaurant(ctx context.Context, ownerID string, req *models.UpdateRestaurantRequest) (*models.Restaurant, error) // ✅ Changed to ownerID
	GetNearbyRestaurants(ctx context.Context, lat, lon, radius float64, search string, isOpenFilter *bool) ([]models.Restaurant, error)
	GetRestaurantByID(ctx context.Context, id string) (*models.Restaurant, error)
}

type restaurantRepository struct {
	db             *mongo.Client
	dbName         string
	collectionName string
}

func NewRestaurantRepository(db *mongo.Client, dbName, collectionName string) RestaurantRepository {
	log.Println("Initializing RestaurantRepository")
	return &restaurantRepository{
		db:             db,
		dbName:         dbName,
		collectionName: collectionName,
	}
}

func (r *restaurantRepository) CheckIfOwnerHasRestaurant(ctx context.Context, ownerID string) (string, bool, error) {
	log.Println("CheckIfOwnerHasRestaurant called with ownerID:", ownerID)

	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	var result models.Restaurant

	err := coll.FindOne(ctx, bson.D{{Key: "owner_id", Value: ownerID}}).Decode(&result) // ✅ Query by owner_id
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("No restaurant found for ownerID:", ownerID)
			return "", false, nil
		}
		log.Println("Error in CheckIfOwnerHasRestaurant:", err)
		return "", true, apperr.ErrInternalServer
	}

	log.Println("Restaurant exists with ID:", result.ID.Hex())
	return result.ID.Hex(), true, nil
}

func (r *restaurantRepository) CreateRestaurant(ctx context.Context, restaurant *models.Restaurant) (*models.Restaurant, error) {
	log.Println("CreateRestaurant called for:", restaurant.Name)

	restaurant.CreatedAt = time.Now()
	restaurant.UpdatedAt = time.Now()

	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	result, err := coll.InsertOne(ctx, restaurant)
	if err != nil {
		log.Println("Error inserting restaurant:", err)
		return nil, apperr.ErrInternalServer
	}

	restaurant.ID = result.InsertedID.(bson.ObjectID)
	log.Println("Restaurant created with ID:", restaurant.ID.Hex())

	return restaurant, nil
}

func (r *restaurantRepository) GetRestaurant(ctx context.Context, ownerID string) (*models.Restaurant, error) {
	log.Println("GetRestaurant called with ownerID:", ownerID)

	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	var result models.Restaurant

	err := coll.FindOne(ctx, bson.D{{Key: "owner_id", Value: ownerID}}).Decode(&result) // ✅ Query by owner_id
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Restaurant not found for ownerID:", ownerID)
			return nil, apperr.ErrRestaurantNotFound
		}
		log.Println("Error in GetRestaurant:", err)
		return nil, apperr.ErrInternalServer
	}

	log.Println("Restaurant fetched:", result.ID.Hex())
	return &result, nil
}

func (r *restaurantRepository) UpdateRestaurant(ctx context.Context, ownerID string, req *models.UpdateRestaurantRequest) (*models.Restaurant, error) {
	log.Println("UpdateRestaurant called for ownerID:", ownerID)

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

	log.Println("Update fields:", updateDocs)

	update := bson.M{"$set": updateDocs}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated models.Restaurant
	err := coll.FindOneAndUpdate(ctx, bson.D{{Key: "owner_id", Value: ownerID}}, update, opts).Decode(&updated) // ✅ Query by owner_id
	if err != nil {
		log.Println("Error updating restaurant:", err)
		return nil, apperr.ErrInternalServer
	}

	log.Println("Restaurant updated:", updated.ID.Hex())
	return &updated, nil
}

func (r *restaurantRepository) GetNearbyRestaurants(ctx context.Context, lat, lon, radius float64, search string, isOpenFilter *bool) ([]models.Restaurant, error) {
	log.Println("GetNearbyRestaurants called with lat:", lat, "lon:", lon, "radius:", radius, "search:", search, "isOpenFilter:", isOpenFilter)

	coll := r.db.Database(r.dbName).Collection(r.collectionName)

	query := bson.M{"is_verified": true}

	if search != "" {
		query["name"] = bson.M{"$regex": search, "$options": "i"}
	}

	if isOpenFilter != nil {
		query["is_open"] = *isOpenFilter
	}

	log.Println("Geo query:", query)

	pipeline := mongo.Pipeline{
		{{Key: "$geoNear", Value: bson.M{
			"near": bson.M{
				"type":        "Point",
				"coordinates": []float64{lon, lat},
			},
			"distanceField": "distanceRaw",
			"maxDistance":   radius,
			"spherical":     true,
			"query":         query,
		}}},
		{{Key: "$sort", Value: bson.D{
			{Key: "is_open", Value: -1},
			{Key: "distanceRaw", Value: 1},
		}}},
		{{Key: "$addFields", Value: bson.M{
			"distanceKm": bson.M{
				"$round": bson.A{
					bson.M{"$divide": bson.A{"$distanceRaw", 1000}},
					2,
				},
			},
		}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Error in aggregation:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var restaurants []models.Restaurant
	if err := cursor.All(ctx, &restaurants); err != nil {
		log.Println("Error decoding restaurants:", err)
		return nil, err
	}

	if restaurants == nil {
		restaurants = []models.Restaurant{}
	}

	log.Println("Nearby restaurants count:", len(restaurants))
	return restaurants, nil
}

func (r *restaurantRepository) GetRestaurantByID(ctx context.Context, id string) (*models.Restaurant, error) {
	log.Println("GetRestaurantByID called with id:", id)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid ObjectID:", err)
		return nil, apperr.ErrInvalidID
	}

	coll := r.db.Database(r.dbName).Collection(r.collectionName)

	var restaurant models.Restaurant
	err = coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&restaurant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Restaurant not found for ID:", id)
			return nil, apperr.ErrRestaurantNotFound
		}
		log.Println("Error fetching restaurant by ID:", err)
		return nil, apperr.ErrInternalServer
	}

	log.Println("Restaurant fetched by ID:", restaurant.ID.Hex())
	return &restaurant, nil
}
