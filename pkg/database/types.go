package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Databaser interface {
	Find(ctx context.Context, filter interface{}, options *options.FindOptions) (*mongo.Cursor, error)
	InsertOne(ctx context.Context, document interface{}) (primitive.ObjectID, *mongo.WriteException)
	Aggregate(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, *mongo.WriteException)
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, *mongo.WriteException)
	Disconnect(ctx context.Context) error
}

// MockDatabase is used to mock the database
type MockDatabase struct {
	FindFunc       func(ctx context.Context, filter interface{}, options *options.FindOptions) (*mongo.Cursor, error)
	InsertOneFunc  func(ctx context.Context, document interface{}) (primitive.ObjectID, *mongo.WriteException)
	AggregateFunc  func(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error)
	UpdateOneFunc  func(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, *mongo.WriteException)
	DeleteOneFunc  func(ctx context.Context, filter interface{}) (*mongo.DeleteResult, *mongo.WriteException)
	DisconnectFunc func(ctx context.Context) error
}

func (s *MockDatabase) Find(ctx context.Context, filter interface{}, options *options.FindOptions) (*mongo.Cursor, error) {
	return s.FindFunc(ctx, filter, options)
}

func (s *MockDatabase) InsertOne(ctx context.Context, document interface{}) (primitive.ObjectID, *mongo.WriteException) {
	return s.InsertOneFunc(ctx, document)
}
func (s *MockDatabase) Aggregate(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
	return s.AggregateFunc(ctx, pipeline)
}

func (s *MockDatabase) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, *mongo.WriteException) {
	return s.UpdateOneFunc(ctx, filter, update)
}

func (s *MockDatabase) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, *mongo.WriteException) {
	return s.DeleteOneFunc(ctx, filter)
}

func (s *MockDatabase) Disconnect(ctx context.Context) error {
	return s.DisconnectFunc(ctx)
}
