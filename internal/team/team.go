package team

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
	"github.com/MorhafAlshibly/coanda/pkg/metric"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	api.UnimplementedTeamServiceServer
	sql                  *sql.DB
	database             *model.Queries
	cache                cache.Cacher
	metric               metric.Metric
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

func WithMetric(metric metric.Metric) func(*Service) {
	return func(input *Service) {
		input.metric = metric
	}
}

func WithMaxMembers(maxMembers uint8) func(*Service) {
	if maxMembers == 0 {
		panic("maxMembers must be greater than 0")
	}
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

func (s *Service) CreateTeam(ctx context.Context, in *api.CreateTeamRequest) (*api.CreateTeamResponse, error) {
	command := NewCreateTeamCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetTeam(ctx context.Context, in *api.GetTeamRequest) (*api.GetTeamResponse, error) {
	command := NewGetTeamCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric).SetInvoker(invoker.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetTeams(ctx context.Context, in *api.GetTeamsRequest) (*api.GetTeamsResponse, error) {
	command := NewGetTeamsCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric).SetInvoker(invoker.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetTeamMember(ctx context.Context, in *api.TeamMemberRequest) (*api.GetTeamMemberResponse, error) {
	command := NewGetTeamMemberCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric).SetInvoker(invoker.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) SearchTeams(ctx context.Context, in *api.SearchTeamsRequest) (*api.SearchTeamsResponse, error) {
	command := NewSearchTeamsCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric).SetInvoker(invoker.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateTeam(ctx context.Context, in *api.UpdateTeamRequest) (*api.UpdateTeamResponse, error) {
	command := NewUpdateTeamCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateTeamMember(ctx context.Context, in *api.UpdateTeamMemberRequest) (*api.UpdateTeamMemberResponse, error) {
	command := NewUpdateTeamMemberCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) DeleteTeam(ctx context.Context, in *api.TeamRequest) (*api.TeamResponse, error) {
	command := NewDeleteTeamCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) JoinTeam(ctx context.Context, in *api.JoinTeamRequest) (*api.JoinTeamResponse, error) {
	command := NewJoinTeamCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) LeaveTeam(ctx context.Context, in *api.TeamMemberRequest) (*api.LeaveTeamResponse, error) {
	command := NewLeaveTeamCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

// Utility functions

func unmarshalTeam(team model.RankedTeam) (*api.Team, error) {
	// Marshal data to protobuf struct
	data, err := conversion.RawJsonToProtobufStruct(team.Data)
	if err != nil {
		return nil, err
	}
	return &api.Team{
		Id:        team.ID,
		Name:      team.Name,
		Score:     team.Score,
		Ranking:   team.Ranking,
		Data:      data,
		CreatedAt: timestamppb.New(team.CreatedAt),
		UpdatedAt: timestamppb.New(team.UpdatedAt),
	}, nil
}

func unmarshalTeamWithMembers(team []model.RankedTeamWithMember) (*api.Team, error) {
	// Marshal data to protobuf struct
	data, err := conversion.RawJsonToProtobufStruct(team[0].Data)
	if err != nil {
		return nil, err
	}
	members := make([]*api.TeamMember, len(team))
	for i, member := range team {
		if !member.MemberID.Valid {
			continue
		}
		memberData, err := conversion.RawJsonToProtobufStruct(member.MemberData)
		if err != nil {
			return nil, err
		}
		members[i] = &api.TeamMember{
			Id:        uint64(member.MemberID.Int64),
			UserId:    uint64(member.UserID.Int64),
			TeamId:    member.ID,
			Data:      memberData,
			JoinedAt:  timestamppb.New(member.JoinedAt.Time),
			UpdatedAt: timestamppb.New(member.MemberUpdatedAt.Time),
		}
	}
	return &api.Team{
		Id:        team[0].ID,
		Name:      team[0].Name,
		Score:     team[0].Score,
		Ranking:   team[0].Ranking,
		Members:   members,
		Data:      data,
		CreatedAt: timestamppb.New(team[0].CreatedAt),
		UpdatedAt: timestamppb.New(team[0].UpdatedAt),
	}, nil
}

func unmarshalTeamMember(member model.TeamMember) (*api.TeamMember, error) {
	// Marshal data to protobuf struct
	data, err := conversion.RawJsonToProtobufStruct(member.Data)
	if err != nil {
		return nil, err
	}
	return &api.TeamMember{
		Id:        member.ID,
		UserId:    member.UserID,
		TeamId:    member.TeamID,
		Data:      data,
		JoinedAt:  timestamppb.New(member.JoinedAt),
		UpdatedAt: timestamppb.New(member.UpdatedAt),
	}, nil
}

func unmarshalTeamsWithMembers(teams []model.RankedTeamWithMember) ([]*api.Team, error) {
	result := make([]*api.Team, 0)
	var currentTeam *api.Team
	for _, team := range teams {
		if currentTeam == nil || currentTeam.Id != team.ID {
			if currentTeam != nil {
				result = append(result, currentTeam)
			}
			data, err := conversion.RawJsonToProtobufStruct(team.Data)
			if err != nil {
				return nil, err
			}
			currentTeam = &api.Team{
				Id:        team.ID,
				Name:      team.Name,
				Score:     team.Score,
				Ranking:   team.Ranking,
				Data:      data,
				CreatedAt: timestamppb.New(team.CreatedAt),
				UpdatedAt: timestamppb.New(team.UpdatedAt),
			}
		}
		if !team.MemberID.Valid {
			continue
		}
		memberData, err := conversion.RawJsonToProtobufStruct(team.MemberData)
		if err != nil {
			return nil, err
		}
		currentTeam.Members = append(currentTeam.Members, &api.TeamMember{
			Id:        uint64(team.MemberID.Int64),
			UserId:    uint64(team.UserID.Int64),
			TeamId:    team.ID,
			Data:      memberData,
			JoinedAt:  timestamppb.New(team.JoinedAt.Time),
			UpdatedAt: timestamppb.New(team.MemberUpdatedAt.Time),
		})
	}
	if currentTeam != nil {
		result = append(result, currentTeam)
	}
	return result, nil
}

// Enum for errors
type TeamRequestError string

const (
	NAME_TOO_SHORT     TeamRequestError = "NAME_TOO_SHORT"
	NAME_TOO_LONG      TeamRequestError = "NAME_TOO_LONG"
	NO_FIELD_SPECIFIED TeamRequestError = "NO_FIELD_SPECIFIED"
)

// Check for errors in team member request
func (s *Service) checkForTeamMemberRequestError(request *api.TeamMemberRequest) *TeamRequestError {
	if request == nil {
		return conversion.ValueToPointer(NO_FIELD_SPECIFIED)
	}
	// Check if id is provided
	if request.Id != nil {
		return nil
	}
	// Check if user id is provided
	if request.UserId != nil {
		return nil
	}
	return conversion.ValueToPointer(NO_FIELD_SPECIFIED)
}

func (s *Service) checkForTeamRequestError(request *api.TeamRequest) *TeamRequestError {
	if request == nil {
		return conversion.ValueToPointer(NO_FIELD_SPECIFIED)
	}
	if request.Id != nil {
		return nil
	}
	// Check if team name is provided
	if request.Name != nil {
		if len(*request.Name) < int(s.minTeamNameLength) {
			return conversion.ValueToPointer(NAME_TOO_SHORT)
		}
		if len(*request.Name) > int(s.maxTeamNameLength) {
			return conversion.ValueToPointer(NAME_TOO_LONG)
		}
		return nil
	}
	return s.checkForTeamMemberRequestError(request.Member)
}
