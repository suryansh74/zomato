package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/suryansh74/zomato/services/auth-service/apperr"
	"github.com/suryansh74/zomato/services/auth-service/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (bson.M, error)
	Create(ctx context.Context, LoginRequest *models.LoginRequest) (*mongo.InsertOneResult, error)
}

type authRepository struct {
	db             *mongo.Client
	dbName         string
	collectionName string
}

func NewAuthRepository(db *mongo.Client, dbName, collectionName string) AuthRepository {
	return &authRepository{
		db:             db,
		dbName:         dbName,
		collectionName: collectionName,
	}
}

func (r *authRepository) FindByEmail(ctx context.Context, email string) (bson.M, error) {
	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	var result bson.M
	err := coll.FindOne(context.Background(), bson.D{{Key: "email", Value: email}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperr.ErrUserNotFound
		}
		return nil, apperr.ErrInternalServer
	}
	return result, nil
}

func (r *authRepository) Create(ctx context.Context, LoginRequest *models.LoginRequest) (*mongo.InsertOneResult, error) {
	newUser := &models.User{
		Name:  LoginRequest.Name,
		Email: LoginRequest.Email,
		Image: LoginRequest.Image,
		// Role:      "customer",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	result, err := coll.InsertOne(context.TODO(), newUser)
	if err != nil {
		return nil, apperr.ErrInternalServer
	}

	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)

	return result, nil
}
