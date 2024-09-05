package task

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/task/model"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
	"github.com/MorhafAlshibly/coanda/pkg/metric"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	api.UnimplementedTaskServiceServer
	sql                  *sql.DB
	database             *model.Queries
	cache                cache.Cacher
	metric               metric.Metric
	defaultMaxPageLength uint8
	maxMaxPageLength     uint8
}

func WithSql(sql *sql.DB) func(*Service) {
	return func(input *Service) {
		input.sql = sql
	}
}

func WithDatabase(database *model.Queries) func(*Service) {
	return func(input *Service) {
		input.database = database
	}
}

func WithCache(cache cache.Cacher) func(*Service) {
	return func(input *Service) {
		input.cache = cache
	}
}

func WithMetric(metric metric.Metric) func(*Service) {
	return func(input *Service) {
		input.metric = metric
	}
}

func WithDefaultMaxPageLength(defaultMaxPageLength uint8) func(*Service) {
	return func(input *Service) {
		input.defaultMaxPageLength = defaultMaxPageLength
	}
}

func WithMaxMaxPageLength(maxMaxPageLength uint8) func(*Service) {
	return func(input *Service) {
		input.maxMaxPageLength = maxMaxPageLength
	}
}

func NewService(opts ...func(*Service)) *Service {
	service := Service{
		defaultMaxPageLength: 10,
		maxMaxPageLength:     100,
	}
	for _, opt := range opts {
		opt(&service)
	}
	return &service
}

func (s *Service) CreateTask(ctx context.Context, input *api.CreateTaskRequest) (*api.CreateTaskResponse, error) {
	command := NewCreateTaskCommand(s, input)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetTask(ctx context.Context, input *api.TaskRequest) (*api.GetTaskResponse, error) {
	command := NewGetTaskCommand(s, input)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric).SetInvoker(invoker.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetTasks(ctx context.Context, input *api.GetTasksRequest) (*api.GetTasksResponse, error) {
	command := NewGetTasksCommand(s, input)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric).SetInvoker(invoker.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateTask(ctx context.Context, input *api.UpdateTaskRequest) (*api.UpdateTaskResponse, error) {
	command := NewUpdateTaskCommand(s, input)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) CompleteTask(ctx context.Context, input *api.TaskRequest) (*api.CompleteTaskResponse, error) {
	command := NewCompleteTaskCommand(s, input)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) DeleteTask(ctx context.Context, input *api.TaskRequest) (*api.TaskResponse, error) {
	command := NewDeleteTaskCommand(s, input)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func unmarshalTask(task *model.Task) (*api.Task, error) {
	data, err := conversion.RawJsonToProtobufStruct(task.Data)
	if err != nil {
		return nil, err
	}
	return &api.Task{
		Id:          task.ID,
		Type:        task.Type,
		Data:        data,
		ExpiresAt:   timestamppb.New(task.ExpiresAt.Time),
		CompletedAt: timestamppb.New(task.CompletedAt.Time),
		CreatedAt:   timestamppb.New(task.CreatedAt),
		UpdatedAt:   timestamppb.New(task.UpdatedAt),
	}, nil

}

// Enum for errors
type TaskRequestError string

const (
	ID_REQUIRED   TaskRequestError = "ID_REQUIRED"
	TYPE_REQUIRED TaskRequestError = "TYPE_REQUIRED"
)

func (s *Service) checkForTaskRequestError(request *api.TaskRequest) *TaskRequestError {
	if request == nil {
		return conversion.ValueToPointer(ID_REQUIRED)
	}
	if request.Id == "" {
		return conversion.ValueToPointer(ID_REQUIRED)
	}
	if request.Type == "" {
		return conversion.ValueToPointer(TYPE_REQUIRED)
	}
	return nil
}
