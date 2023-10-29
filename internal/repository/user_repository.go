package repository

import (
	"context"
	"errors"
	"time"

	common_error "github.com/Bit-Bridge-Source/BitBridge-CommonService-Go/public/error"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// IUserRepository defines the interface for user repository operations.
type IUserRepository interface {
	FindById(ctx context.Context, id string) (*model.PrivateUserModel, error)
	FindByEmail(ctx context.Context, email string) (*model.PrivateUserModel, error)
	FindByUsername(ctx context.Context, username string) (*model.PrivateUserModel, error)
	Create(ctx context.Context, user *model.PrivateUserModel) error
	Update(ctx context.Context, user *model.PrivateUserModel) error
	Delete(ctx context.Context, user *model.PrivateUserModel) error
}

// MongoUserRepository is an implementation of IUserRepository using MongoDB.
type MongoUserRepository struct {
	Collection IUserMongoAdapter
}

// NewUserRepository creates a new instance of MongoUserRepository.
func NewUserRepository(collection IUserMongoAdapter) *MongoUserRepository {
	return &MongoUserRepository{
		Collection: collection,
	}
}

// Create implements IUserRepository.
func (m *MongoUserRepository) Create(ctx context.Context, user *model.PrivateUserModel) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	_, err := m.Collection.InsertOne(ctx, user)

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return common_error.NewServiceError(common_error.Conflict, "User creation failed", err)
		}
	}

	return err
}

// Delete implements IUserRepository.
func (m *MongoUserRepository) Delete(ctx context.Context, user *model.PrivateUserModel) error {
	_, err := m.Collection.DeleteOne(ctx, user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return common_error.NewServiceError(common_error.NotFound, "User not found", err)
		}
	}

	return err
}

// FindByEmail implements IUserRepository.
func (m *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*model.PrivateUserModel, error) {
	var user model.PrivateUserModel
	filter := bson.M{"email": email}
	err := m.Collection.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, common_error.NewServiceError(common_error.NotFound, "User not found", err)
		} else {
			return nil, err
		}
	}

	return &user, err
}

// FindById implements IUserRepository.
func (m *MongoUserRepository) FindById(ctx context.Context, id string) (*model.PrivateUserModel, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user model.PrivateUserModel
	err = m.Collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, common_error.NewServiceError(common_error.NotFound, "User not found", err)
		} else {
			return nil, err
		}
	}

	return &user, err
}

// FindByUsername implements IUserRepository.
func (m *MongoUserRepository) FindByUsername(ctx context.Context, username string) (*model.PrivateUserModel, error) {
	var user model.PrivateUserModel
	filter := bson.M{"username": username}
	err := m.Collection.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, common_error.NewServiceError(common_error.NotFound, "User not found", err)
		} else {
			return nil, err
		}
	}

	return &user, err
}

// Update implements IUserRepository.
func (m *MongoUserRepository) Update(ctx context.Context, user *model.PrivateUserModel) error {
	user.UpdatedAt = time.Now()
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}
	_, err := m.Collection.UpdateOne(ctx, filter, update)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return common_error.NewServiceError(common_error.NotFound, "User not found", err)
		} else {
			return err
		}
	}

	return err
}

var _ IUserRepository = (*MongoUserRepository)(nil)
