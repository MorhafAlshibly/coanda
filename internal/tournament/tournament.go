package tournament

import (
	"context"
	"database/sql"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/tournament/model"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"google.golang.org/protobuf/types/known/timestamppb"

	// "github.com/MorhafAlshibly/coanda/pkg/database"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
)

type Service struct {
	api.UnimplementedTournamentServiceServer
	sql                     *sql.DB
	database                *model.Queries
	cache                   cache.Cacher
	metrics                 metrics.Metrics
	minTournamentNameLength uint8
	maxTournamentNameLength uint8
	dailyTournamentMinute   uint16
	weeklyTournamentMinute  uint16
	weeklyTournamentDay     time.Weekday
	monthlyTournamentMinute uint16
	monthlyTournamentDay    uint8
	defaultMaxPageLength    uint8
	maxMaxPageLength        uint8
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

func WithDailyTournamentMinute(dailyTournamentMinute uint16) func(*Service) {
	return func(input *Service) {
		input.dailyTournamentMinute = dailyTournamentMinute
	}
}

func WithWeeklyTournamentMinute(weeklyTournamentMinute uint16) func(*Service) {
	return func(input *Service) {
		input.weeklyTournamentMinute = weeklyTournamentMinute
	}
}

func WithWeeklyTournamentDay(weeklyTournamentDay time.Weekday) func(*Service) {
	return func(input *Service) {
		input.weeklyTournamentDay = weeklyTournamentDay
	}
}

func WithMonthlyTournamentMinute(monthlyTournamentMinute uint16) func(*Service) {
	return func(input *Service) {
		input.monthlyTournamentMinute = monthlyTournamentMinute
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
		dailyTournamentMinute:   0,
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

func (s *Service) CreateTournamentUser(ctx context.Context, in *api.CreateTournamentUserRequest) (*api.CreateTournamentUserResponse, error) {
	command := NewCreateTournamentUserCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetTournamentUser(ctx context.Context, in *api.TournamentUserRequest) (*api.TournamentUserResponse, error) {
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

func (s *Service) UpdateTournamentUser(ctx context.Context, in *api.UpdateTournamentUserRequest) (*api.UpdateTournamentUserResponse, error) {
	command := NewUpdateTournamentUserCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) DeleteTournamentUser(ctx context.Context, in *api.TournamentUserRequest) (*api.TournamentUserResponse, error) {
	command := NewDeleteTournamentUserCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func unmarshalTournamentUser(tournamentUser *model.RankedTournament) (*api.TournamentUser, error) {
	data, err := conversion.RawJsonToProtobufStruct(tournamentUser.Data)
	if err != nil {
		return nil, err
	}
	var intervalString string
	err = tournamentUser.TournamentInterval.Scan(&intervalString)
	if err != nil {
		return nil, err
	}
	interval := api.TournamentInterval(api.TournamentInterval_value[intervalString])
	return &api.TournamentUser{
		Tournament:          tournamentUser.Name,
		UserId:              tournamentUser.UserID,
		Interval:            interval,
		Score:               tournamentUser.Score,
		Ranking:             tournamentUser.Ranking,
		Data:                data,
		TournamentStartedAt: timestamppb.New(tournamentUser.TournamentStartedAt),
		CreatedAt:           timestamppb.New(tournamentUser.CreatedAt),
		UpdatedAt:           timestamppb.New(tournamentUser.UpdatedAt),
	}, nil
}

func (s *Service) getTournamentStartDate(currentTime time.Time, interval api.TournamentInterval) time.Time {
	var startDate time.Time
	switch interval {
	case api.TournamentInterval_DAILY:
		startDate = currentTime.UTC().Truncate(time.Hour * 24).Add(time.Duration(s.dailyTournamentMinute) * time.Minute)
		if currentTime.UTC().Before(startDate) {
			startDate = startDate.Add(-24 * time.Hour)
		}
	case api.TournamentInterval_WEEKLY:
		startDate = currentTime.UTC().Truncate(time.Hour * 24).Add(time.Duration((int(s.weeklyTournamentDay)-int(currentTime.UTC().Weekday())-7)%7) * 24 * time.Hour).Add(time.Duration(s.weeklyTournamentMinute) * time.Minute)
		if currentTime.UTC().Before(startDate) {
			startDate = startDate.Add(-7 * 24 * time.Hour)
		}
	case api.TournamentInterval_MONTHLY:
		startDate = time.Date(currentTime.UTC().Year(), currentTime.UTC().Month(), int(s.monthlyTournamentDay), 0, 0, 0, 0, time.UTC).Add(time.Duration(s.monthlyTournamentMinute) * time.Minute)
		if currentTime.UTC().Before(startDate) {
			startDate = startDate.AddDate(0, -1, 0)
		}
	default:
		startDate = time.Unix(0, 0).UTC()
	}
	return startDate
}

// Enum for errors
type tournamentUserRequestError string

const (
	NOT_FOUND                 tournamentUserRequestError = "NOT_FOUND"
	TOURNAMENT_NAME_TOO_SHORT tournamentUserRequestError = "TOURNAMENT_NAME_TOO_SHORT"
	TOURNAMENT_NAME_TOO_LONG  tournamentUserRequestError = "TOURNAMENT_NAME_TOO_LONG"
	USER_ID_REQUIRED          tournamentUserRequestError = "USER_ID_REQUIRED"
)

func (s *Service) checkForTournamentUserRequestError(request *api.TournamentUserRequest) *tournamentUserRequestError {
	if request == nil {
		return conversion.ValueToPointer(NOT_FOUND)
	}
	if len(request.Tournament) < int(s.minTournamentNameLength) {
		return conversion.ValueToPointer(TOURNAMENT_NAME_TOO_SHORT)
	}
	if len(request.Tournament) > int(s.maxTournamentNameLength) {
		return conversion.ValueToPointer(TOURNAMENT_NAME_TOO_LONG)
	}
	if request.UserId == 0 {
		return conversion.ValueToPointer(USER_ID_REQUIRED)
	}
	return nil
}

func validateANullTournamentInterval(interval *api.TournamentInterval) model.NullTournamentTournamentInterval {
	if interval == nil {
		return model.NullTournamentTournamentInterval{
			Valid:                        false,
			TournamentTournamentInterval: model.TournamentTournamentIntervalDaily,
		}
	}
	return model.NullTournamentTournamentInterval{Valid: true, TournamentTournamentInterval: model.TournamentTournamentInterval(interval.String())}
}
