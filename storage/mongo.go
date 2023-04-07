//go:build integration
// +build integration

package storage

import (
	"context"
	"fmt"

	errors "github.com/LeoCBS/httpmiddleware/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (mc *MongoClient) InsertRecord(
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
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (mc *MongoClient) DropDatabase(
	ctx context.Context,
	database string,
) error {
	err := mc.client.Database(database).Drop(ctx)

	if err != nil {
		return fmt.Errorf("error to drop database / err = %w", err)
	}
	return nil
}

func (mc *MongoClient) FindByID(
	ctx context.Context,
	database,
	collection string,
	_id string,
) (map[string]interface{}, error) {
	coll := mc.client.Database(database).Collection(collection)
	objectID, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error on create objectID / err = {%v}", err))
	}

	filter := bson.D{{"_id", objectID}}

	var result map[string]interface{}
	err = coll.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFound(fmt.Sprintf("id {%s} not found", _id))
		}
		return nil, errors.New(fmt.Sprintf("error on storage.FindByID / err = {%v}", err))
	}
	return result, nil
}
