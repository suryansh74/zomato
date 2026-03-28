package repositories

import (
	"context"
	"log"
	"time"

	"github.com/suryansh74/zomato/services/auth-service/apperr"
	"github.com/suryansh74/zomato/services/auth-service/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
	Create(ctx context.Context, user *models.User) (*models.User, error)
	UpdateRole(ctx context.Context, role models.Role, id string) (*models.User, error)
}

type authRepository struct {
	db             *mongo.Client
	dbName         string
	collectionName string
}

func NewAuthRepository(db *mongo.Client, dbName, collectionName string) AuthRepository {
	log.Println("[AuthRepository] Initializing repository")
	return &authRepository{
		db:             db,
		dbName:         dbName,
		collectionName: collectionName,
	}
}

func (r *authRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	log.Printf("[FindByEmail] Searching user with email: %s\n", email)

	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	var result models.User

	err := coll.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("[FindByEmail] No user found with email: %s\n", email)
			return nil, apperr.ErrUserNotFound
		}
		log.Printf("[FindByEmail] Error fetching user with email %s: %v\n", email, err)
		return nil, apperr.ErrInternalServer
	}

	log.Printf("[FindByEmail] User found: %s\n", result.ID.Hex())
	return &result, nil
}

func (r *authRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	log.Printf("[FindByID] Searching user with ID: %s\n", id)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("[FindByID] Invalid ObjectID: %s, error: %v\n", id, err)
		return nil, apperr.ErrInternalServer
	}

	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	var result models.User

	err = coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("[FindByID] No user found with ID: %s\n", id)
			return nil, apperr.ErrUserNotFound
		}
		log.Printf("[FindByID] Error fetching user with ID %s: %v\n", id, err)
		return nil, apperr.ErrInternalServer
	}

	log.Printf("[FindByID] User found: %s\n", result.ID.Hex())
	return &result, nil
}

func (r *authRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	log.Printf("[Create] Creating user with email: %s\n", user.Email)

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	res, err := coll.InsertOne(ctx, user)
	if err != nil {
		log.Printf("[Create] Error creating user with email %s: %v\n", user.Email, err)
		return nil, apperr.ErrInternalServer
	}

	// ✅ FIX: Assign the newly generated MongoDB ID back to the user struct!
	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		user.ID = oid
		log.Printf("[Create] User created with ID: %s\n", user.ID.Hex())
	} else {
		log.Printf("[Create] Warning: Failed to parse InsertedID to ObjectID")
	}

	return user, nil
}

func (r *authRepository) UpdateRole(ctx context.Context, role models.Role, id string) (*models.User, error) {
	log.Printf("[UpdateRole] Updating role for user ID: %s to role: %s\n", id, role.Role)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("[UpdateRole] Invalid ObjectID: %s, error: %v\n", id, err)
		return nil, apperr.ErrInternalServer
	}

	coll := r.db.Database(r.dbName).Collection(r.collectionName)

	_, err = coll.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"role": role.Role, "updated_at": time.Now()}},
	)
	if err != nil {
		log.Printf("[UpdateRole] Error updating role for user ID %s: %v\n", id, err)
		return nil, apperr.ErrInternalServer
	}

	log.Printf("[UpdateRole] Role updated successfully for user ID: %s\n", id)

	return r.FindByID(ctx, id)
}
