package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

type MongoDatabaseInput struct {
	connection string
	database   string
	collection string
	indices    []mongo.IndexModel
}

func NewMongoDatabase(ctx context.Context, input MongoDatabaseInput) (*MongoDatabase, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(input.connection))
	if err != nil {
		return nil, err
	}
	database := client.Database(input.database)
	collection := database.Collection(input.collection)
	for _, index := range input.indices {
		_, err := collection.Indexes().CreateOne(ctx, index)
		if err != nil {
			return nil, err
		}
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

func (d *MongoDatabase) FindOne(ctx context.Context, document *interface{}, filter bson.D, pipeline mongo.Pipeline) error {
	cursor, err := d.collection.Aggregate(ctx, pipeline)
	defer cursor.Close(ctx)
	if err != nil {
		return err
	}
	cursor.FindOne(ctx, filter).Decode(document)
	return nil
}
