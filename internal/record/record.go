package record

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/record/model"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/MorhafAlshibly/coanda/pkg/invoker"
	"github.com/MorhafAlshibly/coanda/pkg/metric"
)

type Service struct {
	api.UnimplementedRecordServiceServer
	sql                  *sql.DB
	database             *model.Queries
	cache                cache.Cacher
	metric               metric.Metric
	minRecordNameLength  uint8
	maxRecordNameLength  uint8
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
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetRecord(ctx context.Context, in *api.RecordRequest) (*api.GetRecordResponse, error) {
	command := NewGetRecordCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric).SetInvoker(invoker.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetRecords(ctx context.Context, in *api.GetRecordsRequest) (*api.GetRecordsResponse, error) {
	command := NewGetRecordsCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric).SetInvoker(invoker.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateRecord(ctx context.Context, in *api.UpdateRecordRequest) (*api.UpdateRecordResponse, error) {
	command := NewUpdateRecordCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) DeleteRecord(ctx context.Context, in *api.RecordRequest) (*api.DeleteRecordResponse, error) {
	command := NewDeleteRecordCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func unmarshalRecord(record *model.RankedRecord) (*api.Record, error) {
	data, err := conversion.RawJsonToProtobufStruct(record.Data)
	if err != nil {
		return nil, err
	}
	return &api.Record{
		Id:        record.ID,
		Name:      record.Name,
		UserId:    record.UserID,
		Record:    record.Record,
		Data:      data,
		Ranking:   record.Ranking,
		CreatedAt: timestamppb.New(record.CreatedAt),
		UpdatedAt: timestamppb.New(record.UpdatedAt),
	}, nil
}

func convertNameUserIdToNullNameUserId(nameUserId *api.NameUserId) model.NullNameUserId {
	if nameUserId == nil {
		return model.NullNameUserId{
			Valid: false,
		}
	}
	return model.NullNameUserId{
		Name:   nameUserId.Name,
		UserId: int64(nameUserId.UserId),
		Valid:  true,
	}
}

// Enum for errors
type RecordRequestError string

const (
	ID_OR_NAME_USER_ID_REQUIRED RecordRequestError = "ID_OR_NAME_USER_ID_REQUIRED"
	NAME_TOO_SHORT              RecordRequestError = "NAME_TOO_SHORT"
	NAME_TOO_LONG               RecordRequestError = "NAME_TOO_LONG"
	USER_ID_REQUIRED            RecordRequestError = "USER_ID_REQUIRED"
)

func (s *Service) checkForRecordRequestError(request *api.RecordRequest) *RecordRequestError {
	if request == nil {
		return conversion.ValueToPointer(ID_OR_NAME_USER_ID_REQUIRED)
	}
	if request.Id != nil {
		return nil
	}
	if request.NameUserId == nil {
		return conversion.ValueToPointer(ID_OR_NAME_USER_ID_REQUIRED)
	}
	if len(request.NameUserId.Name) < int(s.minRecordNameLength) {
		return conversion.ValueToPointer(NAME_TOO_SHORT)
	}
	if len(request.NameUserId.Name) > int(s.maxRecordNameLength) {
		return conversion.ValueToPointer(NAME_TOO_LONG)
	}
	if request.NameUserId.UserId == 0 {
		return conversion.ValueToPointer(USER_ID_REQUIRED)
	}
	return nil
}
