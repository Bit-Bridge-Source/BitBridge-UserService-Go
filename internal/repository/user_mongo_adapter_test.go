package repository_test

import (
	"context"
	"testing"

	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockMongoAdapter struct {
	mock.Mock
}

func (m *MockMongoAdapter) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockMongoAdapter) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockMongoAdapter) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func (m *MockMongoAdapter) FindOne(ctx context.Context, filter interface{}) repository.ISingleResult {
	args := m.Called(ctx, filter)
	return args.Get(0).(repository.ISingleResult)
}

func TestDeleteOne(t *testing.T) {
	// Create an instance of the mock adapter
	mockAdapter := new(MockMongoAdapter)

	// Define the expected behavior of the mock DeleteOne method
	mockAdapter.On("DeleteOne", mock.Anything, mock.Anything).Return(&mongo.DeleteResult{}, nil)

	// Create an instance of the UserMongoAdapter with the mock adapter
	userAdapter := repository.NewMongoAdapter(mockAdapter)

	// Call the DeleteOne method in your adapter
	ctx := context.TODO()
	filter := bson.M{"_id": "some_id"} // Replace with your filter
	deleteResult, err := userAdapter.DeleteOne(ctx, filter)

	// Assertions
	assert.Nil(t, err)
	assert.NotNil(t, deleteResult) // Check that the result is not nil

	// Assert that the mock DeleteOne method was called with the expected arguments
	mockAdapter.AssertCalled(t, "DeleteOne", ctx, filter)
}

func TestDeleteOne_Success(t *testing.T) {
	// Create an instance of the mock adapter
	mockAdapter := new(MockMongoAdapter)

	// Define the expected behavior of the mock DeleteOne method
	mockAdapter.On("DeleteOne", mock.Anything, mock.Anything).Return(&mongo.DeleteResult{DeletedCount: 1}, nil)

	// Create an instance of the UserMongoAdapter with the mock adapter
	userAdapter := repository.NewMongoAdapter(mockAdapter)

	// Call the DeleteOne method in your adapter
	ctx := context.TODO()
	filter := bson.M{"_id": "some_id"} // Replace with your filter
	deleteResult, err := userAdapter.DeleteOne(ctx, filter)

	// Assertions
	assert.Nil(t, err)
	assert.NotNil(t, deleteResult)                       // Check that the result is not nil
	assert.Equal(t, int64(1), deleteResult.DeletedCount) // Check the deleted count

	// Assert that the mock DeleteOne method was called with the expected arguments
	mockAdapter.AssertCalled(t, "DeleteOne", ctx, filter)
}

func TestInsertOne_Success(t *testing.T) {
	// Create an instance of the mock adapter
	mockAdapter := new(MockMongoAdapter)

	// Define the expected behavior of the mock InsertOne method
	mockAdapter.On("InsertOne", mock.Anything, mock.Anything).Return(&mongo.InsertOneResult{InsertedID: "some_id"}, nil)

	// Create an instance of the UserMongoAdapter with the mock adapter
	userAdapter := repository.NewMongoAdapter(mockAdapter)

	// Call the InsertOne method in your adapter
	ctx := context.TODO()
	document := bson.M{"key": "value"} // Replace with your document
	insertResult, err := userAdapter.InsertOne(ctx, document)

	// Assertions
	assert.Nil(t, err)
	assert.NotNil(t, insertResult)                      // Check that the result is not nil
	assert.Equal(t, "some_id", insertResult.InsertedID) // Check the inserted ID

	// Assert that the mock InsertOne method was called with the expected arguments
	mockAdapter.AssertCalled(t, "InsertOne", ctx, document)
}

func TestUpdateOne_Success(t *testing.T) {
	// Create an instance of the mock adapter
	mockAdapter := new(MockMongoAdapter)

	// Define the expected behavior of the mock UpdateOne method
	mockAdapter.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil)

	// Create an instance of the UserMongoAdapter with the mock adapter
	userAdapter := repository.NewMongoAdapter(mockAdapter)

	// Call the UpdateOne method in your adapter
	ctx := context.TODO()
	filter := bson.M{"_id": "some_id"}               // Replace with your filter
	update := bson.M{"$set": bson.M{"key": "value"}} // Replace with your update
	updateResult, err := userAdapter.UpdateOne(ctx, filter, update)

	// Assertions
	assert.Nil(t, err)
	assert.NotNil(t, updateResult)                        // Check that the result is not nil
	assert.Equal(t, int64(1), updateResult.MatchedCount)  // Check the matched count
	assert.Equal(t, int64(1), updateResult.ModifiedCount) // Check the modified count

	// Assert that the mock UpdateOne method was called with the expected arguments
	mockAdapter.AssertCalled(t, "UpdateOne", ctx, filter, update)
}

func TestFindOne_Success(t *testing.T) {
	// Create an instance of the mock adapter
	mockAdapter := new(MockMongoAdapter)

	// Define the expected behavior of the mock FindOne method
	mockResult := &mongo.SingleResult{} // Replace with a mock result
	mockAdapter.On("FindOne", mock.Anything, mock.Anything).Return(mockResult)

	// Create an instance of the UserMongoAdapter with the mock adapter
	userAdapter := repository.NewMongoAdapter(mockAdapter)

	// Call the FindOne method in your adapter
	ctx := context.TODO()
	filter := bson.M{"_id": "some_id"} // Replace with your filter
	singleResult := userAdapter.FindOne(ctx, filter)

	// Assertions
	assert.NotNil(t, singleResult) // Check that the result is not nil

	// Additional assertions as needed
}
