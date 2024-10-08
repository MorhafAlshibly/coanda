package item

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/item/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/errorcode"
	"github.com/go-sql-driver/mysql"
)

type CreateItemCommand struct {
	service *Service
	In      *api.CreateItemRequest
	Out     *api.CreateItemResponse
}

func NewCreateItemCommand(service *Service, in *api.CreateItemRequest) *CreateItemCommand {
	return &CreateItemCommand{
		service: service,
		In:      in,
	}
}

func (c *CreateItemCommand) Execute(ctx context.Context) error {
	if c.In.Id == "" {
		c.Out = &api.CreateItemResponse{
			Success: false,
			Error:   api.CreateItemResponse_ID_REQUIRED,
		}
		return nil
	}
	if c.In.Type == "" {
		c.Out = &api.CreateItemResponse{
			Success: false,
			Error:   api.CreateItemResponse_TYPE_REQUIRED,
		}
		return nil
	}
	if c.In.Data == nil {
		c.Out = &api.CreateItemResponse{
			Success: false,
			Error:   api.CreateItemResponse_DATA_REQUIRED,
		}
		return nil
	}
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	_, err = c.service.database.CreateItem(ctx, model.CreateItemParams{
		ID:        c.In.Id,
		Type:      c.In.Type,
		Data:      data,
		ExpiresAt: conversion.TimestampToSqlNullTime(c.In.ExpiresAt),
	})
	// Check if the item already exists
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == errorcode.MySQLErrorCodeDuplicateEntry {
			c.Out = &api.CreateItemResponse{
				Success: false,
				Error:   api.CreateItemResponse_ALREADY_EXISTS,
			}
			return nil
		}
		return err
	}
	c.Out = &api.CreateItemResponse{
		Success: true,
		Error:   api.CreateItemResponse_NONE,
	}
	return nil
}
