package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// MockDatabase is used to mock the database
type MockDatabase struct {
	InsertOneFunc  func(ctx context.Context, document interface{}) (string, error)
	AggregateFunc  func(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error)
	UpdateOneFunc  func(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	DeleteOneFunc  func(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
	DisconnectFunc func(ctx context.Context) error
}

func (s *MockDatabase) InsertOne(ctx context.Context, document interface{}) (string, error) {
	return s.InsertOneFunc(ctx, document)
}
func (s *MockDatabase) Aggregate(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
	return s.AggregateFunc(ctx, pipeline)
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
