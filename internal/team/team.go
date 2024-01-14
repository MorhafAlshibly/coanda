package team

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/database/sqlc"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
)

type Service struct {
	api.UnimplementedTeamServiceServer
	sql                  *sql.DB
	database             *sqlc.Queries
	cache                cache.Cacher
	metrics              metrics.Metrics
	maxMembers           uint8
	minTeamNameLength    uint8
	maxTeamNameLength    uint8
	defaultMaxPageLength uint8
	maxMaxPageLength     uint8
}

func WithSql(sql *sql.DB) func(*Service) {
	return func(input *Service) {
		input.sql = sql
	}
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
	return s.database.Disconnect(ctx)
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

func (s *Service) UpdateTeam(ctx context.Context, in *api.UpdateTeamRequest) (*api.TeamResponse, error) {
	command := NewUpdateTeamCommand(s, in)
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
