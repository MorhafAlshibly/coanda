package storage

import (
	"context"

	"github.com/MorhafAlshibly/coanda/pkg"
	"github.com/MorhafAlshibly/coanda/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseAdapterStorage struct {
	database database.Databaser
}

func NewDatabaseAdapterStorage(database database.Databaser) *DatabaseAdapterStorage {
	return &DatabaseAdapterStorage{
		database: database,
	}
}

func (s *DatabaseAdapterStorage) Add(ctx context.Context, pk string, data map[string]string) (*Object, error) {
	id, writeErr := s.database.InsertOne(ctx, bson.D{
		{Key: "pk", Value: pk},
		{Key: "data", Value: data},
	})
	if writeErr != nil {
		return nil, writeErr
	}
	return &Object{
		Key:  id.Hex(),
		Pk:   pk,
		Data: data,
	}, nil
}

func (s *DatabaseAdapterStorage) Get(ctx context.Context, key string, pk string) (*Object, error) {
	cursor, err := s.database.Find(ctx, bson.D{
		{Key: "_id", Value: key},
		{Key: "pk", Value: pk},
	}, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		var result map[string]interface{}
		decodeErr := cursor.Decode(&result)
		if decodeErr != nil {
			return nil, decodeErr
		}
		return &Object{
			Key:  key,
			Pk:   pk,
			Data: result["data"].(map[string]string),
		}, nil
	}
	return nil, &ObjectNotFoundError{}
}

func (s *DatabaseAdapterStorage) Query(ctx context.Context, filter map[string]any, max int32, page int) ([]*Object, error) {
	int64Max := int64(max)
	int64Page := int64(page)
	cursor, err := s.database.Find(ctx, pkg.MapStringAnyToBsonD(filter), &options.FindOptions{
		Limit: &int64Max,
		Skip:  &int64Page,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []*Object
	for cursor.Next(ctx) {
		var result map[string]interface{}
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, &Object{
			Key:  result["_id"].(primitive.ObjectID).Hex(),
			Pk:   result["pk"].(string),
			Data: result["data"].(map[string]string),
		})
	}
	return results, &PageNotFoundError{}
}
