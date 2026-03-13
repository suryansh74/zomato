package repositories

import (
	"context"
	"time"

	"github.com/suryansh74/zomato/services/auth-service/apperr"
	"github.com/suryansh74/zomato/services/auth-service/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) (*models.User, error)
	UpdateRole(ctx context.Context, role models.Role, email string) (*models.User, error) // ✅ returns *models.User now
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

func (r *authRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	var result models.User // ✅ decode into User struct directly
	err := coll.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperr.ErrUserNotFound
		}
		return nil, apperr.ErrInternalServer
	}
	return &result, nil
}

func (r *authRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	_, err := coll.InsertOne(ctx, user)
	if err != nil {
		return nil, apperr.ErrInternalServer
	}
	return user, nil
}

func (r *authRepository) UpdateRole(ctx context.Context, role models.Role, email string) (*models.User, error) {
	coll := r.db.Database(r.dbName).Collection(r.collectionName)
	_, err := coll.UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": bson.M{"role": role.Role, "updated_at": time.Now()}})
	if err != nil {
		return nil, apperr.ErrInternalServer
	}
	// fetch updated user to return with new role
	return r.FindByEmail(ctx, email)
}
