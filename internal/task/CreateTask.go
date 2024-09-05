package task

import (
	"context"
	"errors"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/task/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/errorcode"
	"github.com/go-sql-driver/mysql"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateTaskCommand struct {
	service *Service
	In      *api.CreateTaskRequest
	Out     *api.CreateTaskResponse
}

func NewCreateTaskCommand(service *Service, in *api.CreateTaskRequest) *CreateTaskCommand {
	return &CreateTaskCommand{
		service: service,
		In:      in,
	}
}

func (c *CreateTaskCommand) Execute(ctx context.Context) error {
	if c.In.Id == "" {
		c.Out = &api.CreateTaskResponse{
			Success: false,
			Error:   api.CreateTaskResponse_ID_REQUIRED,
		}
		return nil
	}
	if c.In.Type == "" {
		c.Out = &api.CreateTaskResponse{
			Success: false,
			Error:   api.CreateTaskResponse_TYPE_REQUIRED,
		}
		return nil
	}
	if c.In.Data == nil {
		c.Out = &api.CreateTaskResponse{
			Success: false,
			Error:   api.CreateTaskResponse_DATA_REQUIRED,
		}
		return nil
	}
	if c.In.ExpiresAt == nil {
		// If the task has no expiry, set it to the Unix epoch time (0)
		c.In.ExpiresAt = timestamppb.New(time.Unix(0, 0))
	}
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	_, err = c.service.database.CreateTask(ctx, model.CreateTaskParams{
		ID:        c.In.Id,
		Type:      c.In.Type,
		Data:      data,
		ExpiresAt: conversion.TimestampToSqlNullTime(c.In.ExpiresAt),
	})
	// Check if the task already exists
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == errorcode.MySQLErrorCodeDuplicateEntry {
			c.Out = &api.CreateTaskResponse{
				Success: false,
				Error:   api.CreateTaskResponse_ALREADY_EXISTS,
			}
			return nil
		}
		return err
	}
	c.Out = &api.CreateTaskResponse{
		Success: true,
		Error:   api.CreateTaskResponse_NONE,
	}
	return nil
}
