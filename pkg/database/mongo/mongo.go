package mongo

import (
	"context"

	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	cache      cache.Cacher
}

type MongoDatabaseInput struct {
	Connection string
	Database   string
	Collection string
	Indices    []mongo.IndexModel
}

func NewMongoDatabase(ctx context.Context, input MongoDatabaseInput, cache cache.Cacher) (*MongoDatabase, error) {
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
		cache:      cache,
	}, nil
}

func (d *MongoDatabase) Disconnect(ctx context.Context) error {
	return d.client.Disconnect(ctx)
}

func (d *MongoDatabase) Find(ctx context.Context, filter interface{}, options *options.FindOptions) (*[]bson.M, error) {
	cmd := NewFindCommand(d, &FindInput{
		Filter:  filter,
		Options: options,
	})
	invoker := invokers.NewCacheInvoker(d.cache)
	err := invoker.Invoke(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return cmd.Out.Result, nil
}

func (d *MongoDatabase) InsertOne(ctx context.Context, document interface{}) (primitive.ObjectID, error) {
	cmd := NewInsertOneCommand(d, &InsertOneInput{
		Document: document,
	})
	invoker := invokers.NewBasicInvoker()
	err := invoker.Invoke(ctx, cmd)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return cmd.Out.InsertedID, nil
}

func (d *MongoDatabase) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	cmd := NewUpdateOneCommand(d, &UpdateOneInput{
		Filter: filter,
		Update: update,
	})
	invoker := invokers.NewBasicInvoker()
	err := invoker.Invoke(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return cmd.Out.Result, nil
}

func (d *MongoDatabase) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	cmd := NewDeleteOneCommand(d, &DeleteOneInput{
		Filter: filter,
	})
	invoker := invokers.NewBasicInvoker()
	err := invoker.Invoke(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return cmd.Out.Result, nil
}
