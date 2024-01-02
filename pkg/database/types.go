package database

import (
	"context"

	"github.com/MorhafAlshibly/coanda/pkg/database/dynamoTable"
)

type Databaser interface {
	PutItem(ctx context.Context, input *dynamoTable.PutItemInput) error
	GetItem(ctx context.Context, input *dynamoTable.GetItemInput) (map[string]any, error)
	UpdateItem(ctx context.Context, input *dynamoTable.UpdateItemInput) error
	DeleteItem(ctx context.Context, input *dynamoTable.DeleteItemInput) error
	Query(ctx context.Context, input *dynamoTable.QueryInput) ([]map[string]any, error)
	Scan(ctx context.Context, input *dynamoTable.ScanInput) ([]map[string]any, error)
}

// MockDatabase is used to mock the database
type MockDatabase struct {
	PutItemFunc    func(ctx context.Context, input *dynamoTable.PutItemInput) error
	GetItemFunc    func(ctx context.Context, input *dynamoTable.GetItemInput) (map[string]any, error)
	UpdateItemFunc func(ctx context.Context, input *dynamoTable.UpdateItemInput) error
	DeleteItemFunc func(ctx context.Context, input *dynamoTable.DeleteItemInput) error
	QueryFunc      func(ctx context.Context, input *dynamoTable.QueryInput) ([]map[string]any, error)
	ScanFunc       func(ctx context.Context, input *dynamoTable.ScanInput) ([]map[string]any, error)
}

func (m *MockDatabase) PutItem(ctx context.Context, input *dynamoTable.PutItemInput) error {
	return m.PutItemFunc(ctx, input)
}

func (m *MockDatabase) GetItem(ctx context.Context, input *dynamoTable.GetItemInput) (map[string]any, error) {
	return m.GetItemFunc(ctx, input)
}

func (m *MockDatabase) UpdateItem(ctx context.Context, input *dynamoTable.UpdateItemInput) error {
	return m.UpdateItemFunc(ctx, input)
}

func (m *MockDatabase) DeleteItem(ctx context.Context, input *dynamoTable.DeleteItemInput) error {
	return m.DeleteItemFunc(ctx, input)
}

func (m *MockDatabase) Query(ctx context.Context, input *dynamoTable.QueryInput) ([]map[string]any, error) {
	return m.QueryFunc(ctx, input)
}

func (m *MockDatabase) Scan(ctx context.Context, input *dynamoTable.ScanInput) ([]map[string]any, error) {
	return m.ScanFunc(ctx, input)
}
