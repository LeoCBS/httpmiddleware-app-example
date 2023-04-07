//go:build integration
// +build integration

package storage

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	client *mongo.Client
}

func NewMongoClient(URI string) (*MongoClient, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	if err != nil {
		return nil, fmt.Errorf("error to create mongo client / err = %w", err)
	}
	return &MongoClient{
		client: client,
	}, nil
}

func (mc *MongoClient) Insert(
	ctx context.Context,
	database,
	collection string,
	record interface{},
) (string, error) {
	coll := mc.client.Database(database).Collection(collection)

	result, err := coll.InsertOne(ctx, record)
	if err != nil {
		return "", fmt.Errorf("error to insert record / err = %w", err)
	}
	return result.InsertedID.(string), nil
}
