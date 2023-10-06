package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/model"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUserRepository interface {
	FindById(ctx context.Context, id string) (*model.PrivateUserModel, error)
	FindByEmail(ctx context.Context, email string) (*model.PrivateUserModel, error)
	FindByUsername(ctx context.Context, username string) (*model.PrivateUserModel, error)
	Create(ctx context.Context, user *model.PrivateUserModel) error
	Update(ctx context.Context, user *model.PrivateUserModel) error
	Delete(ctx context.Context, user *model.PrivateUserModel) error
}

type MongoUserRepository struct {
	Collection IUserMongoAdapter
}

// Create implements IUserRepository.
func (m *MongoUserRepository) Create(ctx context.Context, user *model.PrivateUserModel) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	_, err := m.Collection.InsertOne(ctx, user)

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fiber.NewError(fiber.StatusConflict, "User already exists")
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
	}

	return err
}

// Delete implements IUserRepository.
func (m *MongoUserRepository) Delete(ctx context.Context, user *model.PrivateUserModel) error {
	_, err := m.Collection.DeleteOne(ctx, user)
	return err
}

// FindByEmail implements IUserRepository.
func (m *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*model.PrivateUserModel, error) {
	var user model.PrivateUserModel
	filter := bson.M{"email": email}
	err := m.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return nil, err
	}
	return &user, nil
}

// FindById implements IUserRepository.
func (m *MongoUserRepository) FindById(ctx context.Context, id string) (*model.PrivateUserModel, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}
	var user model.PrivateUserModel
	err = m.Collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername implements IUserRepository.
func (m *MongoUserRepository) FindByUsername(ctx context.Context, username string) (*model.PrivateUserModel, error) {
	var user model.PrivateUserModel
	filter := bson.M{"username": username}
	err := m.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return nil, err
	}
	return &user, nil
}

// Update implements IUserRepository.
func (m *MongoUserRepository) Update(ctx context.Context, user *model.PrivateUserModel) error {
	user.UpdatedAt = time.Now()
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}
	_, err := m.Collection.UpdateOne(ctx, filter, update)
	return err
}

var _ IUserRepository = (*MongoUserRepository)(nil)
