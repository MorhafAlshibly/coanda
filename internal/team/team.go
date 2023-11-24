package team

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/database"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	api.UnimplementedTeamServiceServer
	db                   database.Databaser
	cache                cache.Cacher
	metrics              metrics.Metrics
	maxMembers           uint8
	minTeamNameLength    uint8
	maxTeamNameLength    uint8
	defaultMaxPageLength uint8
	maxMaxPageLength     uint8
}

var (
	// Pipeline to sort by score and add rank
	pipeline = mongo.Pipeline{
		{{Key: "$sort", Value: bson.D{{Key: "score", Value: -1}}}},
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: nil}, {Key: "documents", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}},
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

func WithMaxMembers(maxMembers uint8) func(*Service) {
	return func(input *Service) {
		input.maxMembers = maxMembers
	}
}

func WithMinTeamNameLength(minTeamNameLength uint8) func(*Service) {
	return func(input *Service) {
		input.minTeamNameLength = minTeamNameLength
	}
}

func WithMaxTeamNameLength(maxTeamNameLength uint8) func(*Service) {
	return func(input *Service) {
		input.maxTeamNameLength = maxTeamNameLength
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
	s := &Service{
		maxMembers:           10,
		minTeamNameLength:    3,
		maxTeamNameLength:    20,
		defaultMaxPageLength: 10,
		maxMaxPageLength:     100,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Service) Disconnect(ctx context.Context) error {
	return s.db.Disconnect(ctx)
}

func (s *Service) CreateTeam(ctx context.Context, in *api.CreateTeamRequest) (*api.CreateTeamResponse, error) {
	command := NewCreateTeamCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetTeam(ctx context.Context, in *api.GetTeamRequest) (*api.GetTeamResponse, error) {
	command := NewGetTeamCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetTeams(ctx context.Context, in *api.GetTeamsRequest) (*api.GetTeamsResponse, error) {
	command := NewGetTeamsCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) SearchTeams(ctx context.Context, in *api.SearchTeamsRequest) (*api.SearchTeamsResponse, error) {
	command := NewSearchTeamsCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateTeamScore(ctx context.Context, in *api.UpdateTeamScoreRequest) (*api.TeamResponse, error) {
	command := NewUpdateTeamScoreCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateTeamData(ctx context.Context, in *api.UpdateTeamDataRequest) (*api.TeamResponse, error) {
	command := NewUpdateTeamDataCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) DeleteTeam(ctx context.Context, in *api.GetTeamRequest) (*api.TeamResponse, error) {
	command := NewDeleteTeamCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) JoinTeam(ctx context.Context, in *api.JoinTeamRequest) (*api.JoinTeamResponse, error) {
	command := NewJoinTeamCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) LeaveTeam(ctx context.Context, in *api.LeaveTeamRequest) (*api.LeaveTeamResponse, error) {
	command := NewLeaveTeamCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func getFilter(input *api.GetTeamRequest) (bson.D, error) {
	if input.Id != nil {
		id, err := primitive.ObjectIDFromHex(*input.Id)
		if err != nil {
			return nil, err
		}
		return bson.D{
			{Key: "_id", Value: id},
		}, nil
	}
	if input.Name != nil {
		return bson.D{
			{Key: "name", Value: input.Name},
		}, nil
	}
	if input.Owner != nil {
		if *input.Owner != 0 {
			return bson.D{
				{Key: "owner", Value: input.Owner},
			}, nil
		}
	}
	return nil, errors.New("Invalid input")
}

func toTeam(cursor *mongo.Cursor) (*api.Team, error) {
	var result *bson.M
	err := cursor.Decode(&result)
	if err != nil {
		return nil, err
	}
	// Convert []int64 to []uint64
	membersWithoutOwner := []uint64{}
	for _, member := range (*result)["membersWithoutOwner"].(primitive.A) {
		membersWithoutOwner = append(membersWithoutOwner, uint64(member.(int64)))
	}
	(*result)["membersWithoutOwner"] = membersWithoutOwner
	// Convert data to map[string]string
	var data primitive.M
	if (*result)["data"] != nil {
		data = (*result)["data"].(primitive.M)
	} else {
		data = primitive.M{}
	}
	(*result)["data"] = map[string]string{}
	for key, value := range data {
		(*result)["data"].(map[string]string)[key] = value.(string)
	}
	// If rank is not given, set it to 0
	if _, ok := (*result)["rank"]; !ok {
		(*result)["rank"] = int32(0)
	}
	return &api.Team{
		Id:                  (*result)["_id"].(primitive.ObjectID).Hex(),
		Name:                (*result)["name"].(string),
		Owner:               pkg.InterfaceToUint64((*result)["owner"]),
		MembersWithoutOwner: membersWithoutOwner,
		Score:               pkg.InterfaceToInt64((*result)["score"]),
		Rank:                pkg.InterfaceToUint64((*result)["rank"]),
		Data:                (*result)["data"].(map[string]string),
	}, nil
}
