package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Databaser interface {
	InsertOne(ctx context.Context, document interface{}) (string, error)
	Aggregate(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error)
	Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
}
