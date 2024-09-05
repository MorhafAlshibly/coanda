package task

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/task/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetTasksCommand struct {
	service *Service
	In      *api.GetTasksRequest
	Out     *api.GetTasksResponse
}

func NewGetTasksCommand(service *Service, in *api.GetTasksRequest) *GetTasksCommand {
	return &GetTasksCommand{
		service: service,
		In:      in,
	}
}

func (c *GetTasksCommand) Execute(ctx context.Context) error {
	limit, offset := conversion.PaginationToLimitOffset(c.In.Pagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	result, err := c.service.database.GetTasks(ctx, model.GetTasksParams{
		Type:      conversion.StringToSqlNullString(c.In.Type),
		Completed: conversion.BoolToSqlNullBool(c.In.Completed),
		Limit:     uint64(limit),
		Offset:    uint64(offset),
	})
	if err != nil {
		return err
	}
	tasks := make([]*api.Task, len(result))
	for i, task := range result {
		tasks[i], err = unmarshalTask(&task)
		if err != nil {
			return err
		}
	}
	c.Out = &api.GetTasksResponse{
		Success: true,
		Tasks:   tasks,
	}
	return nil
}
