package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewDatabase(connStr string) (*Database, error) {
	clientOptions := options.Client().ApplyURI(connStr).SetMaxPoolSize(50)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	db := client.Database("user")
	collection := db.Collection("user")
	return &Database{
		Client:     client,
		Collection: collection,
	}, nil
}

func (d *Database) CreateIndexes() error {
	indexModels := []mongo.IndexModel{
		{Keys: map[string]interface{}{"username": 1}, Options: options.Index().SetUnique(true)},
		{Keys: map[string]interface{}{"email": 1}, Options: options.Index().SetUnique(true)},
	}

	_, err := d.Collection.Indexes().CreateMany(context.Background(), indexModels)
	return err
}

func (d *Database) Disconnect() {
	if err := d.Client.Disconnect(context.Background()); err != nil {
		log.Fatal("Failed to disconnect from MongoDB", err)
	}
}
