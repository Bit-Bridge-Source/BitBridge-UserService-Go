package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/model"
	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockMongoOperations struct {
	mock.Mock
}

func (m *MockMongoOperations) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockMongoOperations) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockMongoOperations) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func (m *MockMongoOperations) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.SingleResult)
}

type SingleResultWrapper struct {
	decoder SingleResultDecoder
}

func (s *SingleResultWrapper) Decode(v interface{}) error {
	return s.decoder.Decode(v)
}

type SingleResultDecoder interface {
	Decode(v interface{}) error
}

type MockSingleResult struct {
	mock.Mock
}

func (m *MockSingleResult) Decode(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	mockMongo := new(MockMongoOperations)
	mockInsertOneResult := &mongo.InsertOneResult{}
	repo := repository.NewUserRepository(mockMongo)

	ctx := context.Background()
	user := &model.PrivateUserModel{
		ID:       primitive.NewObjectID(),
		Username: "test",
		Email:    "test@mail.com",
		Hash:     "test",
	}

	// Mock setup: expect InsertOne() to be called with context and user, return mock result and nil error
	mockMongo.On("InsertOne", ctx, mock.MatchedBy(func(u *model.PrivateUserModel) bool {
		return u != nil && u.Email != "" && u.Username != ""
	})).Return(mockInsertOneResult, nil)

	err := repo.Create(ctx, user)

	// Assertions
	assert.Nil(t, err)
	mockMongo.AssertExpectations(t) // Ensure mock expectations are met

	// Clear mock expectations
	mockMongo.ExpectedCalls = nil
}

func TestDuplicateUser(t *testing.T) {
	mockMongo := new(MockMongoOperations)
	mockInsertOneResult := &mongo.InsertOneResult{}
	repo := repository.NewUserRepository(mockMongo)
	ctx := context.Background()
	user := &model.PrivateUserModel{
		ID:       primitive.NewObjectID(),
		Username: "test",
		Email:    "test@mail.com",
		Hash:     "test",
	}

	mockMongo.On("InsertOne", ctx, mock.Anything).Return(mockInsertOneResult, mongo.WriteException{WriteErrors: []mongo.WriteError{{Code: 11000}}})

	err := repo.Create(ctx, user)
	assert.NotNil(t, err)
	mockMongo.AssertExpectations(t)

	// Clear mock expectations
	mockMongo.ExpectedCalls = nil
}

func TestUpdateUser(t *testing.T) {
	mockMongo := new(MockMongoOperations)
	mockUpdateResult := &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}
	repo := repository.NewUserRepository(mockMongo)

	ctx := context.Background()
	user := &model.PrivateUserModel{
		ID:       primitive.NewObjectID(),
		Username: "test",
		Email:    "test@mail.com",
		Hash:     "test",
	}

	// Mock setup: expect UpdateOne() to be called with context and user, return mock result and nil error
	mockMongo.On("UpdateOne", ctx, bson.M{"_id": user.ID}, mock.MatchedBy(func(u bson.M) bool {
		return u != nil
	})).Return(mockUpdateResult, nil)

	err := repo.Update(ctx, user)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, int64(1), mockUpdateResult.MatchedCount)
	assert.Equal(t, int64(1), mockUpdateResult.ModifiedCount)

	mockMongo.AssertExpectations(t) // Ensure mock expectations are met

	// Clear mock expectations
	mockMongo.ExpectedCalls = nil
}

func TestDeleteUser(t *testing.T) {
	mockMongo := new(MockMongoOperations)
	mockDeleteResult := &mongo.DeleteResult{
		DeletedCount: 1,
	}
	repo := repository.NewUserRepository(mockMongo)

	ctx := context.Background()
	user := &model.PrivateUserModel{
		ID:       primitive.NewObjectID(),
		Username: "test",
		Email:    "test@mail.com",
		Hash:     "test",
	}

	// Mock setup: expect DeleteOne() to be called with context and user, return mock result and nil error
	mockMongo.On("DeleteOne", ctx, mock.MatchedBy(func(u *model.PrivateUserModel) bool {
		return u != nil && u.ID != primitive.NilObjectID
	})).Return(mockDeleteResult, nil)

	err := repo.Delete(ctx, user)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, int64(1), mockDeleteResult.DeletedCount)
	mockMongo.AssertExpectations(t) // Ensure mock expectations are met

	// Clear mock expectations
	mockMongo.ExpectedCalls = nil
}

func TestFindById(t *testing.T) {
	ctx := context.TODO()
	id := primitive.NewObjectID()

	// Mock setup
	mockMongo := new(MockMongoOperations)

	// Setup expectedUser
	expectedUser := &model.PrivateUserModel{
		ID: id,
	}

	// Create a real mongo.SingleResult
	sr := mongo.NewSingleResultFromDocument(expectedUser, nil, bson.DefaultRegistry)

	// Setup mockMongo to return the real mongo.SingleResult
	mockMongo.On("FindOne", ctx, bson.M{"_id": id}).Return(sr)

	// Create userRepo with the mocked MongoDB operations
	userRepo := repository.NewUserRepository(mockMongo)

	// Test the method
	user, err := userRepo.FindById(ctx, id.Hex())

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)

	// Verify that the mock expectations were met
	mockMongo.AssertExpectations(t)
}

func TestFindById_UserNotFound(t *testing.T) {
	ctx := context.TODO()
	id := primitive.NewObjectID()

	mockMongo := new(MockMongoOperations)

	sr := mongo.NewSingleResultFromDocument(&model.PrivateUserModel{}, mongo.ErrNoDocuments, bson.DefaultRegistry)

	// Setup mockMongo to return an empty result
	mockMongo.On("FindOne", ctx, bson.M{"_id": id}).Return(sr)

	// Create userRepo with the mocked MongoDB operations
	userRepo := repository.NewUserRepository(mockMongo)

	// Test the method
	user, err := userRepo.FindById(ctx, id.Hex())

	// Custom check for your fiber error, adapt as needed
	fiberError, ok := err.(*fiber.Error)
	assert.True(t, ok)
	assert.Equal(t, fiber.StatusNotFound, fiberError.Code)

	// Assertions
	assert.Nil(t, user)

	// Verify that the mock expectations were met
	mockMongo.AssertExpectations(t)

	// Clear mock expectations
	mockMongo.ExpectedCalls = nil
}

func TestFindById_InvalidId(t *testing.T) {
	ctx := context.TODO()
	invalidId := "invalid"

	mockMongo := new(MockMongoOperations)

	userRepo := repository.NewUserRepository(mockMongo)

	_, err := userRepo.FindById(ctx, invalidId)

	// Assertions
	assert.NotNil(t, err)

	fiberError, ok := err.(*fiber.Error)
	assert.True(t, ok)
	assert.Equal(t, fiber.StatusBadRequest, fiberError.Code)

	mockMongo.AssertExpectations(t)

	// Clear mock expectations
	mockMongo.ExpectedCalls = nil
}

func TestFindById_UnknownError(t *testing.T) {
	ctx := context.TODO()
	id := primitive.NewObjectID()

	mockMongo := new(MockMongoOperations)

	unknownError := errors.New("unknown error")

	sr := mongo.NewSingleResultFromDocument(&model.PrivateUserModel{}, unknownError, bson.DefaultRegistry)

	// Setup mockMongo to return an empty result
	mockMongo.On("FindOne", ctx, bson.M{"_id": id}).Return(sr)

	userRepo := repository.NewUserRepository(mockMongo)

	_, err := userRepo.FindById(ctx, id.Hex())

	// Assertions
	assert.NotNil(t, err)
	assert.Equal(t, unknownError, err)

	mockMongo.AssertExpectations(t)

	// Clear mock expectations
	mockMongo.ExpectedCalls = nil
}

func TestFindByUsername(t *testing.T) {
	ctx := context.Background()
	username := "testuser"

	mockMongo := new(MockMongoOperations)

	// Setup mockSingleResult
	expectedUser := &model.PrivateUserModel{
		Username: username,
	}

	sr := mongo.NewSingleResultFromDocument(expectedUser, nil, bson.DefaultRegistry)

	// Setup mockMongo to return the real mongo.SingleResult
	mockMongo.On("FindOne", ctx, bson.M{"username": username}).Return(sr)

	repo := repository.NewUserRepository(mockMongo)
	user, err := repo.FindByUsername(ctx, username)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, expectedUser.Username, user.Username)

	mockMongo.AssertExpectations(t)

	// Clear mock expectations
	mockMongo.ExpectedCalls = nil
}

func TestFindByUsername_UserNotFound(t *testing.T) {
	ctx := context.Background()
	username := "testuser"

	mockMongo := new(MockMongoOperations)

	sr := mongo.NewSingleResultFromDocument(&model.PrivateUserModel{}, mongo.ErrNoDocuments, bson.DefaultRegistry)

	// Setup mockMongo to return an empty result
	mockMongo.On("FindOne", ctx, bson.M{"username": username}).Return(sr)

	repo := repository.NewUserRepository(mockMongo)
	user, err := repo.FindByUsername(ctx, username)

	// Assertions
	assert.Nil(t, user)

	fiberError, ok := err.(*fiber.Error)
	assert.True(t, ok)
	assert.Equal(t, fiber.StatusNotFound, fiberError.Code)

	// Verify mock expectations
	mockMongo.AssertExpectations(t)

	// Clear mock expectations
	mockMongo.ExpectedCalls = nil
}

func TestFindByUsername_UnknownError(t *testing.T) {
	ctx := context.Background()
	username := "testuser"

	mockMongo := new(MockMongoOperations)

	unknownError := errors.New("unknown error")
	sr := mongo.NewSingleResultFromDocument(&model.PrivateUserModel{}, unknownError, bson.DefaultRegistry)

	// Setup mockMongo to return a result with an error
	mockMongo.On("FindOne", ctx, bson.M{"username": username}).Return(sr)

	repo := repository.NewUserRepository(mockMongo)
	user, err := repo.FindByUsername(ctx, username)

	// Assertions
	assert.NotNil(t, err)
	assert.Equal(t, unknownError, err)
	assert.Nil(t, user)

	// Verify mock expectations
	mockMongo.AssertExpectations(t)

	// Clear mock expectations
	mockMongo.ExpectedCalls = nil
}

func TestFindByEmail(t *testing.T) {
	ctx := context.Background()
	email := "test@mail.com"

	mockMongo := new(MockMongoOperations)

	expectedUser := &model.PrivateUserModel{
		Email: email,
	}
	sr := mongo.NewSingleResultFromDocument(expectedUser, nil, bson.DefaultRegistry)
	mockMongo.On("FindOne", ctx, bson.M{"email": email}).Return(sr)

	repo := repository.NewUserRepository(mockMongo)
	user, err := repo.FindByEmail(ctx, email)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, expectedUser.Email, user.Email)

	mockMongo.AssertExpectations(t)
	mockMongo.ExpectedCalls = nil
}

func TestFindByEmail_UserNotFound(t *testing.T) {
	ctx := context.Background()
	email := "test@mail.com"

	mockMongo := new(MockMongoOperations)

	sr := mongo.NewSingleResultFromDocument(&model.PrivateUserModel{}, mongo.ErrNoDocuments, bson.DefaultRegistry)
	mockMongo.On("FindOne", ctx, bson.M{"email": email}).Return(sr)

	repo := repository.NewUserRepository(mockMongo)
	user, err := repo.FindByEmail(ctx, email)

	// Assertions
	assert.Nil(t, user)
	fiberError, ok := err.(*fiber.Error)
	assert.True(t, ok)
	assert.Equal(t, fiber.StatusNotFound, fiberError.Code)

	mockMongo.AssertExpectations(t)
	mockMongo.ExpectedCalls = nil
}

func TestFindByEmail_UnknownError(t *testing.T) {
	ctx := context.Background()
	email := "test@mail.com"

	mockMongo := new(MockMongoOperations)

	unknownError := errors.New("unknown error")
	sr := mongo.NewSingleResultFromDocument(&model.PrivateUserModel{}, unknownError, bson.DefaultRegistry)
	mockMongo.On("FindOne", ctx, bson.M{"email": email}).Return(sr)

	repo := repository.NewUserRepository(mockMongo)
	user, err := repo.FindByEmail(ctx, email)

	// Assertions
	assert.NotNil(t, err)
	assert.Equal(t, unknownError, err)
	assert.Nil(t, user)

	mockMongo.AssertExpectations(t)
	mockMongo.ExpectedCalls = nil
}
