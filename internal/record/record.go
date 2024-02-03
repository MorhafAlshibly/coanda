package record

import (
	"context"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"

	// "github.com/MorhafAlshibly/coanda/pkg/database/sqlc"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
)

type Service struct {
	api.UnimplementedRecordServiceServer
	database             *sqlc.Queries
	cache                cache.Cacher
	metrics              metrics.Metrics
	minRecordNameLength  uint8
	maxRecordNameLength  uint8
	defaultMaxPageLength uint8
	maxMaxPageLength     uint8
}

func WithDatabase(database *sqlc.Queries) func(*Service) {
	return func(input *Service) {
		input.database = database
	}
}

func WithCache(cache cache.Cacher) func(*Service) {
	return func(input *Service) {
		input.cache = cache
	}
}

func WithMetrics(metrics metrics.Metrics) func(*Service) {
	return func(input *Service) {
		input.metrics = metrics
	}
}

func WithMinRecordNameLength(minRecordNameLength uint8) func(*Service) {
	return func(input *Service) {
		input.minRecordNameLength = minRecordNameLength
	}
}

func WithMaxRecordNameLength(maxRecordNameLength uint8) func(*Service) {
	return func(input *Service) {
		input.maxRecordNameLength = maxRecordNameLength
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
		minRecordNameLength:  3,
		maxRecordNameLength:  20,
		defaultMaxPageLength: 10,
		maxMaxPageLength:     100,
	}
	for _, opt := range opts {
		opt(&service)
	}
	return &service
}

func (s *Service) CreateRecord(ctx context.Context, in *api.CreateRecordRequest) (*api.CreateRecordResponse, error) {
	command := NewCreateRecordCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetRecord(ctx context.Context, in *api.GetRecordRequest) (*api.GetRecordResponse, error) {
	command := NewGetRecordCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetRecords(ctx context.Context, in *api.GetRecordsRequest) (*api.GetRecordsResponse, error) {
	command := NewGetRecordsCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateRecord(ctx context.Context, in *api.UpdateRecordRequest) (*api.UpdateRecordResponse, error) {
	command := NewUpdateRecordCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) DeleteRecord(ctx context.Context, in *api.GetRecordRequest) (*api.DeleteRecordResponse, error) {
	command := NewDeleteRecordCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func UnmarshalRecord(record *sqlc.RankedRecord) (*api.Record, error) {
	data, err := conversion.RawJsonToProtobufStruct(record.Data)
	if err != nil {
		return nil, err
	}
	return &api.Record{
		Name:      record.Name,
		UserId:    record.UserID,
		Record:    record.Record,
		Data:      data,
		Ranking:   record.Ranking.(uint64),
		CreatedAt: record.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt: record.UpdatedAt.Time.Format(time.RFC3339),
	}, nil
}
