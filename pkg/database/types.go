package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Databaser interface {
	InsertOne(ctx context.Context, document interface{}) (string, *mongo.WriteException)
	Aggregate(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, *mongo.WriteException)
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, *mongo.WriteException)
	Disconnect(ctx context.Context) error
}

// MockDatabase is used to mock the database
type MockDatabase struct {
	InsertOneFunc  func(ctx context.Context, document interface{}) (string, *mongo.WriteException)
	AggregateFunc  func(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error)
	UpdateOneFunc  func(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, *mongo.WriteException)
	DeleteOneFunc  func(ctx context.Context, filter interface{}) (*mongo.DeleteResult, *mongo.WriteException)
	DisconnectFunc func(ctx context.Context) error
}

func (s *MockDatabase) InsertOne(ctx context.Context, document interface{}) (string, *mongo.WriteException) {
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
