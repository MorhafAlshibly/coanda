package tournament

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/tournament/model"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/tournament"
	"google.golang.org/protobuf/types/known/timestamppb"

	// "github.com/MorhafAlshibly/coanda/pkg/database"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
	"github.com/MorhafAlshibly/coanda/pkg/metric"
)

type Service struct {
	api.UnimplementedTournamentServiceServer
	sql                     *sql.DB
	database                *model.Queries
	cache                   cache.Cacher
	metric                  metric.Metric
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

func WithMetric(metric metric.Metric) func(*Service) {
	return func(input *Service) {
		input.metric = metric
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
		weeklyTournamentMinute:  0,
		weeklyTournamentDay:     time.Monday,
		monthlyTournamentMinute: 0,
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
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetTournamentUser(ctx context.Context, in *api.TournamentUserRequest) (*api.GetTournamentUserResponse, error) {
	command := NewGetTournamentUserCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric).SetInvoker(invoker.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetTournamentUsers(ctx context.Context, in *api.GetTournamentUsersRequest) (*api.GetTournamentUsersResponse, error) {
	command := NewGetTournamentUsersCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric).SetInvoker(invoker.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateTournamentUser(ctx context.Context, in *api.UpdateTournamentUserRequest) (*api.UpdateTournamentUserResponse, error) {
	command := NewUpdateTournamentUserCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) DeleteTournamentUser(ctx context.Context, in *api.TournamentUserRequest) (*api.TournamentUserResponse, error) {
	command := NewDeleteTournamentUserCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
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
	interval := api.TournamentInterval(api.TournamentInterval_value[strings.ToUpper(string(tournamentUser.TournamentInterval))])
	return &api.TournamentUser{
		Id:                  tournamentUser.ID,
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

// Enum for errors
type tournamentUserRequestError string

const (
	NOT_FOUND                                  tournamentUserRequestError = "NOT_FOUND"
	ID_OR_TOURNAMENT_INTERVAL_USER_ID_REQUIRED tournamentUserRequestError = "ID_OR_TOURNAMENT_INTERVAL_USER_ID_REQUIRED"
	TOURNAMENT_NAME_TOO_SHORT                  tournamentUserRequestError = "TOURNAMENT_NAME_TOO_SHORT"
	TOURNAMENT_NAME_TOO_LONG                   tournamentUserRequestError = "TOURNAMENT_NAME_TOO_LONG"
	USER_ID_REQUIRED                           tournamentUserRequestError = "USER_ID_REQUIRED"
)

func (s *Service) checkForTournamentUserRequestError(request *api.TournamentUserRequest) *tournamentUserRequestError {
	if request == nil {
		return conversion.ValueToPointer(ID_OR_TOURNAMENT_INTERVAL_USER_ID_REQUIRED)
	}
	if request.Id == nil {
		if request.TournamentIntervalUserId == nil {
			return conversion.ValueToPointer(ID_OR_TOURNAMENT_INTERVAL_USER_ID_REQUIRED)
		}
		if len(request.TournamentIntervalUserId.Tournament) < int(s.minTournamentNameLength) {
			return conversion.ValueToPointer(TOURNAMENT_NAME_TOO_SHORT)
		}
		if len(request.TournamentIntervalUserId.Tournament) > int(s.maxTournamentNameLength) {
			return conversion.ValueToPointer(TOURNAMENT_NAME_TOO_LONG)
		}
		if request.TournamentIntervalUserId.UserId == 0 {
			return conversion.ValueToPointer(USER_ID_REQUIRED)
		}
		return nil
	}
	return nil
}

func (s *Service) convertTournamentIntervalUserIdToNullNameIntervalUserIDStartedAt(nameIntervalUserID *api.TournamentIntervalUserId) model.NullNameIntervalUserIDStartedAt {
	if nameIntervalUserID == nil {
		return model.NullNameIntervalUserIDStartedAt{
			Valid: false,
		}
	}
	return model.NullNameIntervalUserIDStartedAt{
		Name:               nameIntervalUserID.Tournament,
		TournamentInterval: model.TournamentTournamentInterval(nameIntervalUserID.Interval.String()),
		UserID:             nameIntervalUserID.UserId,
		TournamentStartedAt: tournament.GetStartTime(time.Now().UTC(), nameIntervalUserID.Interval, tournament.WipeTimes{
			DailyTournamentMinute:   s.dailyTournamentMinute,
			WeeklyTournamentMinute:  s.weeklyTournamentMinute,
			WeeklyTournamentDay:     s.weeklyTournamentDay,
			MonthlyTournamentMinute: s.monthlyTournamentMinute,
			MonthlyTournamentDay:    s.monthlyTournamentDay,
		}),
		Valid: true,
	}
}
