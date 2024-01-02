package record

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/database"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	api.UnimplementedRecordServiceServer
	db                   database.Databaser
	cache                cache.Cacher
	metrics              metrics.Metrics
	minRecordNameLength  uint8
	maxRecordNameLength  uint8
	defaultMaxPageLength uint8
	maxMaxPageLength     uint8
}

var (
	// Pipeline to partition by name, then sort by record and add rank
	pipeline = mongo.Pipeline{
		{{Key: "$sort", Value: bson.D{{Key: "name", Value: 1}, {Key: "record", Value: 1}}}},
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$name"}, {Key: "documents", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "name", Value: "$_id"},
			{Key: "documents", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$documents"},
					{Key: "as", Value: "doc"},
					{Key: "in", Value: bson.D{
						{Key: "$mergeObjects", Value: bson.A{
							"$$doc",
							bson.D{
								{Key: "rank", Value: bson.D{
									{Key: "$add", Value: bson.A{
										bson.D{{Key: "$indexOfArray", Value: bson.A{"$documents", "$$doc"}}},
										1,
									}},
								}},
							},
						}},
					}},
				}},
			}},
		}}},
		{{Key: "$unwind", Value: "$documents"}},
		{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$documents"}}}},
	}
)

func WithDatabase(db database.Databaser) func(*Service) {
	return func(input *Service) {
		input.db = db
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

func MarshalRecord(record *api.Record) (map[string]any, error) {
	data, err := conversion.ProtobufStructToMap(record.Data)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"name":      record.Name,
		"userId":    record.UserId,
		"record":    record.Record,
		"data":      data,
		"createdAt": record.CreatedAt,
		"updatedAt": record.UpdatedAt,
	}, nil
}

func UnmarshalRecord(record map[string]any) (*api.Record, error) {
	data, err := conversion.MapToProtobufStruct(record["data"].(map[string]any))
	if err != nil {
		return nil, err
	}
	return &api.Record{
		Name:      record["name"].(string),
		UserId:    record["userId"].(uint64),
		Record:    record["record"].(uint64),
		Rank:      record["rank"].(uint64),
		Data:      data,
		CreatedAt: record["createdAt"].(string),
		UpdatedAt: record["updatedAt"].(string),
	}, nil
}
