package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

type MongoDatabaseInput struct {
	Connection string
	Database   string
	Collection string
	Indices    []mongo.IndexModel
}

func NewMongoDatabase(ctx context.Context, input MongoDatabaseInput) (*MongoDatabase, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(input.Connection))
	if err != nil {
		return nil, err
	}
	database := client.Database(input.Database)
	collection := database.Collection(input.Collection)
	_, err = collection.Indexes().CreateMany(ctx, input.Indices)
	if err != nil {
		return nil, err
	}
	return &MongoDatabase{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

func (d *MongoDatabase) Disconnect(ctx context.Context) error {
	return d.client.Disconnect(ctx)
}

func (d *MongoDatabase) InsertOne(ctx context.Context, document interface{}) (string, error) {
	result, err := d.collection.InsertOne(ctx, document)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).String(), nil
}

func (d *MongoDatabase) Aggregate(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func (d *MongoDatabase) Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func (d *MongoDatabase) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	result, err := d.collection.UpdateOne(ctx, filter, update, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *MongoDatabase) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
