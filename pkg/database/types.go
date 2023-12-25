package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Databaser interface {
	Find(ctx context.Context, filter interface{}, options *options.FindOptions) (*[]bson.M, error)
	InsertOne(ctx context.Context, document interface{}) (primitive.ObjectID, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
	Disconnect(ctx context.Context) error
}

// MockDatabase is used to mock the database
type MockDatabase struct {
	FindFunc       func(ctx context.Context, filter interface{}, options *options.FindOptions) (*[]bson.M, error)
	InsertOneFunc  func(ctx context.Context, document interface{}) (primitive.ObjectID, error)
	UpdateOneFunc  func(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	DeleteOneFunc  func(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
	DisconnectFunc func(ctx context.Context) error
}

func (s *MockDatabase) Find(ctx context.Context, filter interface{}, options *options.FindOptions) (*[]bson.M, error) {
	return s.FindFunc(ctx, filter, options)
}

func (s *MockDatabase) InsertOne(ctx context.Context, document interface{}) (primitive.ObjectID, error) {
	return s.InsertOneFunc(ctx, document)
}

func (s *MockDatabase) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return s.UpdateOneFunc(ctx, filter, update)
}

func (s *MockDatabase) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return s.DeleteOneFunc(ctx, filter)
}

func (s *MockDatabase) Disconnect(ctx context.Context) error {
	return s.DisconnectFunc(ctx)
}
