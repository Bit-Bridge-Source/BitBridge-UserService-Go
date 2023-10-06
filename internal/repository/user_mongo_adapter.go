package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoSingleResultAlias mongo.SingleResult

func (s *MongoSingleResultAlias) Decode(v interface{}) error {
	return (*mongo.SingleResult)(s).Decode(v)
}

type DecodeResult interface {
	Decode(v interface{}) error
}

type IUserMongoAdapter interface {
	InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
	FindOne(ctx context.Context, filter interface{}) DecodeResult
}

type UserMongoAdapter struct {
	Collection *mongo.Collection
}

// DeleteOne implements IUserMongoAdapter.
func (m *UserMongoAdapter) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return m.Collection.DeleteOne(ctx, filter)
}

// FindOne implements IUserMongoAdapter.
func (m *UserMongoAdapter) FindOne(ctx context.Context, filter interface{}) DecodeResult {
	return (*MongoSingleResultAlias)(m.Collection.FindOne(ctx, filter))
}

// InsertOne implements IUserMongoAdapter.
func (m *UserMongoAdapter) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return m.Collection.InsertOne(ctx, document)
}

// UpdateOne implements IUserMongoAdapter.
func (m *UserMongoAdapter) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return m.Collection.UpdateOne(ctx, filter, update)
}

func NewMongoAdapter(collection *mongo.Collection) *UserMongoAdapter {
	return &UserMongoAdapter{
		Collection: collection,
	}
}

var _ IUserMongoAdapter = (*UserMongoAdapter)(nil)
