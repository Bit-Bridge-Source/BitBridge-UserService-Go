package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// ISingleResult represents the result of a query.
type ISingleResult interface {
	Decode(v interface{}) error
}

// SingleResult wraps mongo.SingleResult
type SingleResult struct {
	Result *mongo.SingleResult
}

// Decode decodes the single result.
func (r *SingleResult) Decode(v interface{}) error {
	return r.Result.Decode(v)
}

// IUserMongoAdapter interface
type IUserMongoAdapter interface {
	InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
	FindOne(ctx context.Context, filter interface{}) ISingleResult
	// Add other methods if needed
}

// UserMongoAdapter struct
type UserMongoAdapter struct {
	Adapter IUserMongoAdapter
}

// DeleteOne implements IUserMongoAdapter.
func (m *UserMongoAdapter) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return m.Adapter.DeleteOne(ctx, filter)
}

// FindOne implements IUserMongoAdapter.
func (m *UserMongoAdapter) FindOne(ctx context.Context, filter interface{}) ISingleResult {
	return m.Adapter.FindOne(ctx, filter)
}

// InsertOne implements IUserMongoAdapter.
func (m *UserMongoAdapter) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return m.Adapter.InsertOne(ctx, document)
}

// UpdateOne implements IUserMongoAdapter.
func (m *UserMongoAdapter) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return m.Adapter.UpdateOne(ctx, filter, update)
}

// NewMongoAdapter creates a new UserMongoAdapter
func NewMongoAdapter(adapter IUserMongoAdapter) *UserMongoAdapter {
	return &UserMongoAdapter{
		Adapter: adapter,
	}
}

var _ IUserMongoAdapter = (*UserMongoAdapter)(nil)
