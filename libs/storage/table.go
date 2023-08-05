package storage

import (
	"context"
	"errors"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/bytedance/sonic"
	"github.com/google/uuid"
)

// TableStorage is used to store data in azure table storage
type TableStorage struct {
	Client *aztables.Client
}

// NewTableStorage creates a new table storage
func NewTableStorage(ctx context.Context, connection string, tableName string) (*TableStorage, error) {
	// Create the service client
	serviceClient, err := aztables.NewServiceClientFromConnectionString(connection, nil)
	if err != nil {
		return nil, err
	}
	// Create the table client
	client := serviceClient.NewClient(tableName)
	// Create the table
	_, err = client.CreateTable(ctx, nil)
	if err != nil {
		// If the table already exists, ignore the error
		if !strings.Contains(err.Error(), string(aztables.TableAlreadyExists)) {
			return nil, err
		}
	}
	// Return the table storage
	return &TableStorage{Client: client}, nil
}

// Add is used to add data to the table
func (s *TableStorage) Add(ctx context.Context, pk string, data map[string]any) (string, error) {
	// Generate a new key
	key := uuid.New().String()
	// Create the entity
	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: pk,
			RowKey:       key,
		},
		Properties: data,
	}
	// Marshal the entity
	marshalled, err := sonic.Marshal(entity)
	if err != nil {
		return "", err
	}
	// Add the entity to the table
	_, err = s.Client.AddEntity(ctx, marshalled, nil)
	if err != nil {
		return "", err
	}
	return key, nil
}

// Get is used to get data from the table
func (s *TableStorage) Get(ctx context.Context, key string, pk string) (map[string]any, error) {
	// Get the entity from the table
	entityResponse, err := s.Client.GetEntity(ctx, pk, key, nil)
	if err != nil {
		// If the entity is not found, return an error
		if strings.Contains(err.Error(), string(aztables.ResourceNotFound)) {
			return nil, errors.New("Data not found")
		}
		return nil, err
	}
	// Unmarshal the entity
	var entity map[string]any
	err = sonic.Unmarshal(entityResponse.Value, &entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// Query is used to query data from the table
func (s *TableStorage) Query(ctx context.Context, filter string, max int32, page int) ([]QueryResult, error) {
	// Set the options
	options := &aztables.ListEntitiesOptions{
		Filter: &filter,
		Top:    &max,
	}
	// Create the pager
	pager := s.Client.NewListEntitiesPager(options)
	// Iterate through the pages until the page is found
	pageCount := 0
	for pager.More() {
		// Increment the page count
		pageCount++
		// Get the next page
		entities, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		// Page found
		if pageCount == page {
			// Specify the output
			var out []QueryResult
			for _, entity := range entities.Entities {
				// Unmarshal the entity
				var edmEntity aztables.EDMEntity
				err = sonic.Unmarshal(entity, &edmEntity)
				if err != nil {
					return nil, err
				}
				// Append the output
				out = append(out, QueryResult{
					Key:  edmEntity.Entity.RowKey,
					Pk:   edmEntity.Entity.PartitionKey,
					Data: edmEntity.Properties,
				})

			}
			return out, nil
		}
	}
	return nil, errors.New("Page not found")
}
