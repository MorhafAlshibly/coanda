package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.74

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/bff/model"
)

// CreateItem is the resolver for the CreateItem field.
func (r *mutationResolver) CreateItem(ctx context.Context, input model.CreateItemRequest) (*model.CreateItemResponse, error) {
	resp, err := r.itemClient.CreateItem(ctx, &api.CreateItemRequest{
		Id:        input.ID,
		Type:      input.Type,
		Data:      input.Data,
		ExpiresAt: input.ExpiresAt,
	})
	if err != nil {
		return nil, err
	}
	return &model.CreateItemResponse{
		Success: resp.Success,
		Error:   model.CreateItemError(resp.Error.String()),
	}, nil
}

// UpdateItem is the resolver for the UpdateItem field.
func (r *mutationResolver) UpdateItem(ctx context.Context, input model.UpdateItemRequest) (*model.UpdateItemResponse, error) {
	resp, err := r.itemClient.UpdateItem(ctx, &api.UpdateItemRequest{
		Item: &api.ItemRequest{
			Id:   input.Item.ID,
			Type: input.Item.Type,
		},
		Data: input.Data,
	})
	if err != nil {
		return nil, err
	}
	return &model.UpdateItemResponse{
		Success: resp.Success,
		Error:   model.UpdateItemError(resp.Error.String()),
	}, nil
}

// DeleteItem is the resolver for the DeleteItem field.
func (r *mutationResolver) DeleteItem(ctx context.Context, input model.ItemRequest) (*model.ItemResponse, error) {
	resp, err := r.itemClient.DeleteItem(ctx, &api.ItemRequest{
		Id:   input.ID,
		Type: input.Type,
	})
	if err != nil {
		return nil, err
	}
	return &model.ItemResponse{
		Success: resp.Success,
		Error:   model.ItemError(resp.Error.String()),
	}, nil
}

// GetItem is the resolver for the GetItem field.
func (r *queryResolver) GetItem(ctx context.Context, input model.ItemRequest) (*model.GetItemResponse, error) {
	resp, err := r.itemClient.GetItem(ctx, &api.ItemRequest{
		Id:   input.ID,
		Type: input.Type,
	})
	if err != nil {
		return nil, err
	}
	var item *model.Item
	if resp.Item != nil {
		item = &model.Item{
			ID:        resp.Item.Id,
			Type:      resp.Item.Type,
			Data:      resp.Item.Data,
			ExpiresAt: resp.Item.ExpiresAt,
			CreatedAt: resp.Item.CreatedAt,
			UpdatedAt: resp.Item.UpdatedAt,
		}
	}
	return &model.GetItemResponse{
		Success: resp.Success,
		Item:    item,
		Error:   model.GetItemError(resp.Error.String()),
	}, nil
}

// GetItems is the resolver for the GetItems field.
func (r *queryResolver) GetItems(ctx context.Context, input model.GetItemsRequest) (*model.GetItemsResponse, error) {
	if input.Pagination == nil {
		input.Pagination = &model.Pagination{}
	}
	resp, err := r.itemClient.GetItems(ctx, &api.GetItemsRequest{
		Type: input.Type,
		Pagination: &api.Pagination{
			Max:  input.Pagination.Max,
			Page: input.Pagination.Page,
		},
	})
	if err != nil {
		return nil, err
	}
	items := make([]*model.Item, len(resp.Items))
	for i, item := range resp.Items {
		items[i] = &model.Item{
			ID:        item.Id,
			Type:      item.Type,
			Data:      item.Data,
			ExpiresAt: item.ExpiresAt,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}
	return &model.GetItemsResponse{
		Success: resp.Success,
		Items:   items,
	}, nil
}
