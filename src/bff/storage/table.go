package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/bytedance/sonic"
	"github.com/google/uuid"
)

type TableStorage struct {
	Client *aztables.Client
}

func NewTableStorage(ctx context.Context, connection string, tableName string) (*TableStorage, error) {
	serviceClient, err := aztables.NewServiceClientFromConnectionString(connection, nil)
	if err != nil {
		return nil, err
	}
	client := serviceClient.NewClient(tableName)
	_, err = client.CreateTable(ctx, nil)
	if err != nil {
		if !strings.Contains(err.Error(), string(aztables.TableAlreadyExists)) {
			return nil, err
		}
	}
	return &TableStorage{Client: client}, nil
}

func (s *TableStorage) Add(ctx context.Context, pk string, data map[string]any) (string, error) {
	key := uuid.New().String()
	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: pk,
			RowKey:       key,
		},
		Properties: data,
	}
	marshalled, err := sonic.Marshal(entity)
	if err != nil {
		return "", err
	}
	_, err = s.Client.AddEntity(ctx, marshalled, nil)
	if err != nil {
		fmt.Printf("Error adding entity: %s\n", err.Error())
		return "", err
	}
	return key, nil
}

func (s *TableStorage) Get(ctx context.Context, key string, pk string) (map[string]any, error) {
	entityResponse, err := s.Client.GetEntity(ctx, pk, key, nil)
	if err != nil {
		if strings.Contains(err.Error(), string(aztables.ResourceNotFound)) {
			return nil, errors.New("Data not found")
		}
		return nil, err
	}
	var entity map[string]any
	err = sonic.Unmarshal(entityResponse.Value, &entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *TableStorage) Query(ctx context.Context, filter string, max int32, page int) ([]QueryResult, error) {
	options := &aztables.ListEntitiesOptions{
		Filter: &filter,
		Top:    &max,
	}
	pager := s.Client.NewListEntitiesPager(options)
	pageCount := 0
	for pager.More() {
		pageCount++
		entities, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		if pageCount == page {
			var out []QueryResult
			for _, entity := range entities.Entities {
				var edmEntity aztables.EDMEntity
				err = sonic.Unmarshal(entity, &edmEntity)
				if err != nil {
					return nil, err
				}
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
