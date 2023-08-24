package database

import "context"

type Databaser interface {
	InsertOne(ctx context.Context, document interface{}) (interface{}, error)
	InsertMany(ctx context.Context, documents []interface{}) ([]interface{}, error)
	FindOne(ctx context.Context, filter interface{}) (interface{}, error)
	Find(ctx context.Context, filter interface{}) ([]interface{}, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (interface{}, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}) ([]interface{}, error)
	DeleteOne(ctx context.Context, filter interface{}) (interface{}, error)
	DeleteMany(ctx context.Context, filter interface{}) ([]interface{}, error)
}
