package storage

import (
	"context"
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
	// Create the service client, and create the table if it doesn't exist
	serviceClient, err := aztables.NewServiceClientFromConnectionString(connection, nil)
	if err != nil {
		return nil, err
	}
	client := serviceClient.NewClient(tableName)
	_, err = client.CreateTable(ctx, nil)
	if err != nil {
		// If the table already exists, ignore the error
		if !strings.Contains(err.Error(), string(aztables.TableAlreadyExists)) {
			return nil, err
		}
	}
	return &TableStorage{Client: client}, nil
}

// Add is used to add data to the table
func (s *TableStorage) Add(ctx context.Context, pk string, data map[string]string) (*Object, error) {
	// Generate a new key
	key := uuid.New().String()
	// Create the entity
	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			RowKey:       key,
			PartitionKey: pk,
		},
		Properties: *stringMapToAnyMap(&data),
	}
	marshalled, err := sonic.Marshal(entity)
	if err != nil {
		return nil, err
	}
	// Add the entity to the table
	_, err = s.Client.AddEntity(ctx, marshalled, nil)
	if err != nil {
		return nil, err
	}
	return &Object{
		Key:  key,
		Pk:   pk,
		Data: data,
	}, nil
}

// Get is used to get data from the table
func (s *TableStorage) Get(ctx context.Context, key string, pk string) (*Object, error) {
	entityResponse, err := s.Client.GetEntity(ctx, pk, key, nil)
	if err != nil {
		// If the entity is not found, return an error
		if strings.Contains(err.Error(), string(aztables.ResourceNotFound)) {
			return nil, &ObjectNotFoundError{}
		}
		return nil, err
	}
	return entityToObject(&entityResponse.Value)
}

// Query is used to query data from the table
func (s *TableStorage) Query(ctx context.Context, filter string, max int32, page int) ([]*Object, error) {
	// Set the options and create pager
	options := &aztables.ListEntitiesOptions{
		Filter: &filter,
		Top:    &max,
	}
	pager := s.Client.NewListEntitiesPager(options)
	// Iterate through the pages until the page is found
	pageCount := 0
	for pager.More() {
		pageCount++
		entities, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		// Page found
		if pageCount == page {
			var out []*Object
			for _, entity := range entities.Entities {
				// Convert the entity to an object
				object, err := entityToObject(&entity)
				if err != nil {
					return nil, err
				}
				out = append(out, object)
			}
			return out, nil
		}
	}
	return nil, &PageNotFoundError{}
}

// Helper function to convert entity to object
func entityToObject(entity *[]byte) (*Object, error) {
	var edmEntity aztables.EDMEntity
	err := sonic.Unmarshal(*entity, &edmEntity)
	if err != nil {
		return nil, err
	}
	// Convert entity properties to map[string]string
	properties := make(map[string]string)
	for k, v := range edmEntity.Properties {
		properties[k] = v.(string)
	}
	return &Object{
		Key:  edmEntity.Entity.RowKey,
		Pk:   edmEntity.Entity.PartitionKey,
		Data: properties,
	}, nil
}

// Helper function to convert map[string]string to map[string]any
func stringMapToAnyMap(in *map[string]string) *map[string]any {
	out := make(map[string]any)
	for k, v := range *in {
		out[k] = v
	}
	return &out
}
