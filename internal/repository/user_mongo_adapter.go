package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IUserMongoAdapter interface
type IUserMongoAdapter interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	// Add other methods if needed
}

// UserMongoAdapter struct
type UserMongoAdapter struct {
	Adapter IUserMongoAdapter
}

// DeleteOne implements IUserMongoAdapter.
func (m *UserMongoAdapter) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return m.Adapter.DeleteOne(ctx, filter, opts...)
}

// FindOne implements IUserMongoAdapter.
func (m *UserMongoAdapter) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return m.Adapter.FindOne(ctx, filter)
}

// InsertOne implements IUserMongoAdapter.
func (m *UserMongoAdapter) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return m.Adapter.InsertOne(ctx, document, opts...)
}

// UpdateOne implements IUserMongoAdapter.
func (m *UserMongoAdapter) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return m.Adapter.UpdateOne(ctx, filter, update, opts...)
}

// NewMongoAdapter creates a new UserMongoAdapter
func NewMongoAdapter(adapter IUserMongoAdapter) *UserMongoAdapter {
	return &UserMongoAdapter{
		Adapter: adapter,
	}
}

var _ IUserMongoAdapter = (*UserMongoAdapter)(nil)
