package tournament

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
	api.UnimplementedTournamentServiceServer
	db                      database.Databaser
	cache                   cache.Cacher
	metrics                 metrics.Metrics
	minTournamentNameLength uint8
	maxTournamentNameLength uint8
	weeklyTournamentDay     time.Weekday
	monthlyTournamentDay    uint8
	defaultMaxPageLength    uint8
	maxMaxPageLength        uint8
}

var (
	pipeline = mongo.Pipeline{
		bson.D{
			{Key: "$setWindowFields", Value: bson.D{
				{Key: "partitionBy", Value: bson.D{
					{Key: "tournament", Value: "$tournament"},
					{Key: "interval", Value: "$interval"},
					{Key: "tournamentStartDate", Value: "$tournamentStartDate"},
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

func WithMinTournamentNameLength(minTournamentNameLength uint8) func(*Service) {
	return func(input *Service) {
		input.minTournamentNameLength = minTournamentNameLength
	}
}

func WithMaxTournamentNameLength(maxTournamentNameLength uint8) func(*Service) {
	return func(input *Service) {
		input.maxTournamentNameLength = maxTournamentNameLength
	}
}

func WithWeeklyTournamentDay(weeklyTournamentDay time.Weekday) func(*Service) {
	return func(input *Service) {
		input.weeklyTournamentDay = weeklyTournamentDay
	}
}

func WithMonthlyTournamentDay(monthlyTournamentDay uint8) func(*Service) {
	return func(input *Service) {
		input.monthlyTournamentDay = monthlyTournamentDay
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
		minTournamentNameLength: 3,
		maxTournamentNameLength: 20,
		weeklyTournamentDay:     time.Monday,
		monthlyTournamentDay:    1,
		defaultMaxPageLength:    10,
		maxMaxPageLength:        100,
	}
	for _, opt := range opts {
		opt(&service)
	}
	return &service
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

func (s *Service) GetTournamentUser(ctx context.Context, in *api.GetTournamentUserRequest) (*api.TournamentUserResponse, error) {
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

func (s *Service) UpdateTournamentUserScore(ctx context.Context, in *api.UpdateTournamentUserScoreRequest) (*api.TournamentUserResponse, error) {
	command := NewUpdateTournamentUserScoreCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateTournamentUserData(ctx context.Context, in *api.UpdateTournamentUserDataRequest) (*api.TournamentUserResponse, error) {
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
	if input.Id != "" {
		id, err := primitive.ObjectIDFromHex(input.Id)
		if err != nil {
			return nil, err
		}
		return bson.D{
			{Key: "_id", Value: id},
		}, nil
	}
	if input.TournamentIntervalUserId != nil {
		if input.TournamentIntervalUserId.Tournament != "" {
			if input.TournamentIntervalUserId.UserId != 0 {
				return bson.D{
					{Key: "tournament", Value: input.TournamentIntervalUserId.Tournament},
					{Key: "interval", Value: input.TournamentIntervalUserId.Interval},
				}, nil
			}
		}
	}
	return nil, errors.New("Invalid input")
}

func toTournamentUsers(ctx context.Context, cursor *mongo.Cursor, page uint64, max uint8) ([]*api.TournamentUser, error) {
	var result []*api.TournamentUser
	skip := (int(page) - 1) * int(max)
	for i := 0; i < skip; i++ {
		cursor.Next(ctx)
	}
	for i := 0; i < int(max); i++ {
		if !cursor.Next(ctx) {
			break
		}
		tournament, err := toTournamentUser(cursor)
		if err != nil {
			return nil, err
		}
		result = append(result, tournament)
	}
	return result, nil
}

func toTournamentUser(cursor *mongo.Cursor) (*api.TournamentUser, error) {
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
	return &api.TournamentUser{
		Id:                  (*result)["_id"].(primitive.ObjectID).Hex(),
		UserId:              uint64((*result)["userId"].(int64)),
		Tournament:          (*result)["tournament"].(string),
		Interval:            api.TournamentInterval((*result)["interval"].(int32)),
		Score:               int64((*result)["score"].(int64)),
		Rank:                uint64((*result)["rank"].(int32)),
		TournamentStartDate: (*result)["tournament"].(string),
		Data:                (*result)["data"].(map[string]string),
	}, nil
}

func (s *Service) getTournamentStartDate(currentTime time.Time, interval api.TournamentInterval) string {
	switch interval {
	case api.TournamentInterval_DAILY:
		return currentTime.UTC().Truncate(time.Hour * 24).Format(time.RFC3339)
	case api.TournamentInterval_WEEKLY:
		return currentTime.UTC().Truncate(time.Hour * 24).Add(time.Duration((int(s.weeklyTournamentDay)-int(currentTime.UTC().Weekday())-7)%7) * 24 * time.Hour).Format(time.RFC3339)
	case api.TournamentInterval_MONTHLY:
		return time.Date(currentTime.UTC().Year(), currentTime.UTC().Month(), int(s.monthlyTournamentDay), 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	}
	return time.Time{}.UTC().Format(time.RFC3339)
}
