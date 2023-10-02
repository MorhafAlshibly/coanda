package record

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/database"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	api.UnimplementedRecordServiceServer
	db                   database.Databaser
	cache                cache.Cacher
	metrics              metrics.Metrics
	minRecordNameLength  int
	defaultMaxPageLength uint64
}

type NewServiceInput struct {
	Db                   database.Databaser
	Cache                cache.Cacher
	Metrics              metrics.Metrics
	MinRecordNameLength  int
	DefaultMaxPageLength uint64
}

var (
	pipeline = mongo.Pipeline{
		bson.D{
			{Key: "$setWindowFields", Value: bson.D{
				{Key: "sortBy", Value: bson.D{
					{Key: "record", Value: 1},
				}},
				{Key: "output", Value: bson.D{
					{Key: "rank", Value: bson.D{
						{Key: "$rank", Value: bson.D{}},
					}},
				}},
			}},
		},
	}
)

func NewService(ctx context.Context, input NewServiceInput) *Service {
	return &Service{
		db:                   input.Db,
		cache:                input.Cache,
		metrics:              input.Metrics,
		minRecordNameLength:  input.MinRecordNameLength,
		defaultMaxPageLength: input.DefaultMaxPageLength,
	}
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

func (s *Service) DeleteRecord(ctx context.Context, in *api.DeleteRecordRequest) (*api.DeleteRecordResponse, error) {
	command := NewDeleteRecordCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func getFilter(input *api.GetRecordRequest) (bson.D, error) {
	if input.Id != nil {
		id, err := primitive.ObjectIDFromHex(*input.Id)
		if err != nil {
			return nil, err
		}
		return bson.D{
			{Key: "_id", Value: id},
		}, nil
	}
	if input.NameUserId != nil {
		return bson.D{
			{Key: "name", Value: input.NameUserId.Name},
			{Key: "userId", Value: input.NameUserId.UserId},
		}, nil
	}
	return nil, errors.New("Invalid input")
}

func toRecords(ctx context.Context, cursor *mongo.Cursor, page uint64, max uint64) ([]*api.Record, error) {
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
		CreatedAt: timestamppb.New((*result)["_id"].(primitive.ObjectID).Timestamp()),
	}, nil
}
