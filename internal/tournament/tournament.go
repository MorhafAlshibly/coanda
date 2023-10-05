package tournament

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
	api.UnimplementedTournamentServiceServer
	db                      database.Databaser
	cache                   cache.Cacher
	metrics                 metrics.Metrics
	minTournamentNameLength int
	defaultMaxPageLength    uint64
}

type NewServiceInput struct {
	Db                      database.Databaser
	Cache                   cache.Cacher
	Metrics                 metrics.Metrics
	MinTournamentNameLength int
	DefaultMaxPageLength    uint64
}

var (
	pipeline = mongo.Pipeline{
		bson.D{
			{Key: "$setWindowFields", Value: bson.D{
				{Key: "partitionBy", Value: bson.D{
					{Key: "tournament", Value: "$tournament"},
					{Key: "interval", Value: "$interval"},
				}},
			},
			},
			{Key: "sortBy", Value: bson.D{
				{Key: "score", Value: -1},
			}},
			{Key: "output", Value: bson.D{
				{Key: "rank", Value: bson.D{
					{Key: "$rank", Value: bson.D{}},
				}},
			}},
		}}
)

func NewService(ctx context.Context, input NewServiceInput) *Service {
	return &Service{
		db:                      input.Db,
		cache:                   input.Cache,
		metrics:                 input.Metrics,
		minTournamentNameLength: input.MinTournamentNameLength,
		defaultMaxPageLength:    input.DefaultMaxPageLength,
	}
}

func (s *Service) Disconnect(ctx context.Context) error {
	return s.db.Disconnect(ctx)
}

func (s *Service) CreateTournamentUser(ctx context.Context, in *api.CreateTournamentUserRequest) (*api.CreateTournamentUserResponse, error) {
	command := NewCreateTournamentUserCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetTournamentUser(ctx context.Context, in *api.GetTournamentUserRequest) (*api.GetTournamentUserResponse, error) {
	command := NewGetTournamentUserCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetTournamentUsers(ctx context.Context, in *api.GetTournamentUsersRequest) (*api.GetTournamentUsersResponse, error) {
	command := NewGetTournamentUsersCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateTournamentUserScore(ctx context.Context, in *api.UpdateTournamentUserScoreRequest) (*api.UpdateTournamentUserScoreResponse, error) {
	command := NewUpdateTournamentUserScoreCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateTournamentUserData(ctx context.Context, in *api.UpdateTournamentUserDataRequest) (*api.UpdateTournamentUserDataResponse, error) {
	command := NewUpdateTournamentUserDataCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) DeleteTournamentUser(ctx context.Context, in *api.GetTournamentUserRequest) (*api.DeleteTournamentUserResponse, error) {
	command := NewDeleteTournamentUserCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func getFilter(input *api.GetTournamentUserRequest) (bson.D, error) {
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

func toTournaments(ctx context.Context, cursor *mongo.Cursor, page uint64, max uint64) ([]*api.Tournament, error) {
	var result []*api.Tournament
	skip := (int(page) - 1) * int(max)
	for i := 0; i < skip; i++ {
		cursor.Next(ctx)
	}
	for i := 0; i < int(max); i++ {
		if !cursor.Next(ctx) {
			break
		}
		tournament, err := toTournament(cursor)
		if err != nil {
			return nil, err
		}
		result = append(result, tournament)
	}
	return result, nil
}

func toTournament(cursor *mongo.Cursor) (*api.Tournament, error) {
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
	return &api.Tournament{
		Id:         (*result)["_id"].(primitive.ObjectID).Hex(),
		Name:       (*result)["name"].(string),
		UserId:     uint64((*result)["userId"].(int64)),
		Tournament: uint64((*result)["tournament"].(int64)),
		Rank:       uint64((*result)["rank"].(int32)),
		Data:       (*result)["data"].(map[string]string),
		CreatedAt:  timestamppb.New((*result)["_id"].(primitive.ObjectID).Timestamp()),
	}, nil
}
