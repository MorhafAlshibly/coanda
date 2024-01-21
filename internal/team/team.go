package team

import (
	"context"
	"database/sql"
	"reflect"
	"strconv"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	api.UnimplementedTeamServiceServer
	sql                  *sql.DB
	database             *model.Queries
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

// Utility functions

func UnmarshalTeam(team model.RankedTeam) (*api.Team, error) {
	// Marshal ranking to uint64
	ranking, err := strconv.Atoi(reflect.TypeOf(team.Ranking).String())
	if err != nil {
		return nil, err
	}
	// Marshal data to protobuf struct
	data, err := conversion.RawJsonToProtobufStruct(team.Data)
	if err != nil {
		return nil, err
	}
	return &api.Team{
		Name:      team.Name,
		Owner:     team.Owner,
		Score:     team.Score,
		Ranking:   uint64(ranking),
		Data:      data,
		CreatedAt: timestamppb.New(team.CreatedAt),
		UpdatedAt: timestamppb.New(team.UpdatedAt),
	}, nil
}

func UnmarshalTeamMember(member model.TeamMember) (*api.TeamMember, error) {
	// Marshal data to protobuf struct
	data, err := conversion.RawJsonToProtobufStruct(member.Data)
	if err != nil {
		return nil, err
	}
	return &api.TeamMember{
		Team:      member.Team,
		UserId:    member.UserID,
		Data:      data,
		JoinedAt:  timestamppb.New(member.JoinedAt),
		UpdatedAt: timestamppb.New(member.UpdatedAt),
	}, nil
}

// Enum for result
type GetTeamFieldResult string

const (
	NAME               GetTeamFieldResult = "NAME"
	OWNER              GetTeamFieldResult = "OWNER"
	MEMBER             GetTeamFieldResult = "MEMBER"
	NAME_TOO_SHORT     GetTeamFieldResult = "NAME_TOO_SHORT"
	NAME_TOO_LONG      GetTeamFieldResult = "NAME_TOO_LONG"
	NO_FIELD_SPECIFIED GetTeamFieldResult = "NO_FIELD_SPECIFIED"
)

func (s *Service) GetTeamField(request *api.TeamRequest) GetTeamFieldResult {
	// Check if team name is provided
	if request.Name != nil {
		if len(*request.Name) < int(s.minTeamNameLength) {
			return NAME_TOO_SHORT
		}
		if len(*request.Name) > int(s.maxTeamNameLength) {
			return NAME_TOO_LONG
		}
		return NAME
		// Check if owner is provided
	} else if request.Owner != nil {
		return OWNER
		// Check if member is provided
	} else if request.Member != nil {
		return MEMBER
	} else {
		return NO_FIELD_SPECIFIED
	}
}
