package record

import (
	"context"
	"errors"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/database"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (s *Service) Disconnect(ctx context.Context) error {
	return s.db.Disconnect(ctx)
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

func (s *Service) DeleteRecord(ctx context.Context, in *api.GetRecordRequest) (*api.DeleteRecordResponse, error) {
	command := NewDeleteRecordCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func getFilter(input *api.GetRecordRequest) (bson.D, error) {
	if input.Id != "" {
		id, err := primitive.ObjectIDFromHex(input.Id)
		if err != nil {
			return nil, err
		}
		return bson.D{
			{Key: "_id", Value: id},
		}, nil
	}
	if input.NameUserId != nil {
		if input.NameUserId.UserId != 0 {
			return bson.D{
				{Key: "name", Value: input.NameUserId.Name},
				{Key: "userId", Value: input.NameUserId.UserId},
			}, nil
		}
	}
	return nil, errors.New("Invalid input")
}

func toRecords(ctx context.Context, cursor *mongo.Cursor, page uint64, max uint8) ([]*api.Record, error) {
	var result []*api.Record
	skip := (int(page) - 1) * int(max)
	for i := 0; i < skip; i++ {
		cursor.Next(ctx)
	}
	for i := 0; i < int(max); i++ {
		if !cursor.Next(ctx) {
			break
		}
		record, err := toRecord(cursor)
		if err != nil {
			return nil, err
		}
		result = append(result, record)
	}
	return result, nil
}

func toRecord(cursor *mongo.Cursor) (*api.Record, error) {
	var result *bson.M
	err := cursor.Decode(&result)
	if err != nil {
		return nil, err
	}
	// Convert data to map[string]string
	data := (*result)["data"].(primitive.M)
	(*result)["data"] = map[string]string{}
	for key, value := range data {
		(*result)["data"].(map[string]string)[key] = value.(string)
	}
	// If rank is not given, set it to 0
	if _, ok := (*result)["rank"]; !ok {
		(*result)["rank"] = int32(0)
	}
	return &api.Record{
		Id:        (*result)["_id"].(primitive.ObjectID).Hex(),
		Name:      (*result)["name"].(string),
		UserId:    uint64((*result)["userId"].(int64)),
		Record:    uint64((*result)["record"].(int64)),
		Rank:      uint64((*result)["rank"].(int32)),
		Data:      (*result)["data"].(map[string]string),
		CreatedAt: (*result)["_id"].(primitive.ObjectID).Timestamp().Format(time.RFC3339),
	}, nil
}
