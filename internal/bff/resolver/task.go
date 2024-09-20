package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.54

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/bff/model"
)

// CreateTask is the resolver for the CreateTask field.
func (r *mutationResolver) CreateTask(ctx context.Context, input model.CreateTaskRequest) (*model.CreateTaskResponse, error) {
	resp, err := r.taskClient.CreateTask(ctx, &api.CreateTaskRequest{
		Id:        input.ID,
		Type:      input.Type,
		Data:      input.Data,
		ExpiresAt: input.ExpiresAt,
	})
	if err != nil {
		return nil, err
	}
	return &model.CreateTaskResponse{
		Success: resp.Success,
		Error:   model.CreateTaskError(resp.Error.String()),
	}, nil
}

// UpdateTask is the resolver for the UpdateTask field.
func (r *mutationResolver) UpdateTask(ctx context.Context, input model.UpdateTaskRequest) (*model.UpdateTaskResponse, error) {
	if input.Task == nil {
		input.Task = &model.TaskRequest{}
	}
	resp, err := r.taskClient.UpdateTask(ctx, &api.UpdateTaskRequest{
		Task: &api.TaskRequest{
			Id:   input.Task.ID,
			Type: input.Task.Type,
		},
		Data: input.Data,
	})
	if err != nil {
		return nil, err
	}
	return &model.UpdateTaskResponse{
		Success: resp.Success,
		Error:   model.UpdateTaskError(resp.Error.String()),
	}, nil
}

// CompleteTask is the resolver for the CompleteTask field.
func (r *mutationResolver) CompleteTask(ctx context.Context, input model.TaskRequest) (*model.CompleteTaskResponse, error) {
	resp, err := r.taskClient.CompleteTask(ctx, &api.TaskRequest{
		Id:   input.ID,
		Type: input.Type,
	})
	if err != nil {
		return nil, err
	}
	return &model.CompleteTaskResponse{
		Success: resp.Success,
		Error:   model.CompleteTaskError(resp.Error.String()),
	}, nil
}

// DeleteTask is the resolver for the DeleteTask field.
func (r *mutationResolver) DeleteTask(ctx context.Context, input model.TaskRequest) (*model.TaskResponse, error) {
	resp, err := r.taskClient.DeleteTask(ctx, &api.TaskRequest{
		Id:   input.ID,
		Type: input.Type,
	})
	if err != nil {
		return nil, err
	}
	return &model.TaskResponse{
		Success: resp.Success,
		Error:   model.TaskError(resp.Error.String()),
	}, nil
}

// GetTask is the resolver for the GetTask field.
func (r *queryResolver) GetTask(ctx context.Context, input model.TaskRequest) (*model.GetTaskResponse, error) {
	resp, err := r.taskClient.GetTask(ctx, &api.TaskRequest{
		Id:   input.ID,
		Type: input.Type,
	})
	if err != nil {
		return nil, err
	}
	var task *model.Task
	if resp.Task != nil {
		task = &model.Task{
			ID:          resp.Task.Id,
			Type:        resp.Task.Type,
			Data:        resp.Task.Data,
			ExpiresAt:   resp.Task.ExpiresAt,
			CompletedAt: resp.Task.CompletedAt,
			CreatedAt:   resp.Task.CreatedAt,
			UpdatedAt:   resp.Task.UpdatedAt,
		}
	}
	return &model.GetTaskResponse{
		Success: resp.Success,
		Task:    task,
		Error:   model.GetTaskError(resp.Error.String()),
	}, nil
}

// GetTasks is the resolver for the GetTasks field.
func (r *queryResolver) GetTasks(ctx context.Context, input model.GetTasksRequest) (*model.GetTasksResponse, error) {
	if input.Pagination == nil {
		input.Pagination = &model.Pagination{}
	}
	resp, err := r.taskClient.GetTasks(ctx, &api.GetTasksRequest{
		Type:      input.Type,
		Completed: input.Completed,
		Pagination: &api.Pagination{
			Max:  input.Pagination.Max,
			Page: input.Pagination.Page,
		},
	})
	if err != nil {
		return nil, err
	}
	var tasks []*model.Task
	for _, t := range resp.Tasks {
		tasks = append(tasks, &model.Task{
			ID:          t.Id,
			Type:        t.Type,
			Data:        t.Data,
			ExpiresAt:   t.ExpiresAt,
			CompletedAt: t.CompletedAt,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}
	return &model.GetTasksResponse{
		Success: resp.Success,
		Tasks:   tasks,
	}, nil
}
