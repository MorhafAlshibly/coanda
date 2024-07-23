package matchmaking

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/invokers"
	"github.com/MorhafAlshibly/coanda/pkg/metrics"
)

type Service struct {
	api.UnimplementedMatchmakingServiceServer
	sql                  *sql.DB
	database             *model.Queries
	cache                cache.Cacher
	metrics              metrics.Metrics
	minArenaNameLength   uint8
	maxArenaNameLength   uint8
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

func WithMinArenaNameLength(minArenaNameLength uint8) func(*Service) {
	return func(input *Service) {
		input.minArenaNameLength = minArenaNameLength
	}
}

func WithMaxArenaNameLength(maxArenaNameLength uint8) func(*Service) {
	return func(input *Service) {
		input.maxArenaNameLength = maxArenaNameLength
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

func NewService(options ...func(*Service)) *Service {
	service := Service{
		defaultMaxPageLength: 10,
		maxMaxPageLength:     100,
	}
	for _, option := range options {
		option(&service)
	}
	return &service
}

func (s *Service) CreateArena(ctx context.Context, in *api.CreateArenaRequest) (*api.CreateArenaResponse, error) {
	command := NewCreateArenaCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetArena(ctx context.Context, in *api.ArenaRequest) (*api.GetArenaResponse, error) {
	command := NewGetArenaCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetArenas(ctx context.Context, in *api.GetArenasRequest) (*api.GetArenasResponse, error) {
	command := NewGetArenasCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateArena(ctx context.Context, in *api.UpdateArenaRequest) (*api.UpdateArenaResponse, error) {
	command := NewUpdateArenaCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) DeleteArena(ctx context.Context, in *api.ArenaRequest) (*api.ArenaResponse, error) {
	command := NewDeleteArenaCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) CreateMatchmakingUser(ctx context.Context, in *api.CreateMatchmakingUserRequest) (*api.CreateMatchmakingUserResponse, error) {
	command := NewCreateMatchmakingUserCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetMatchmakingUser(ctx context.Context, in *api.MatchmakingUserRequest) (*api.GetMatchmakingUserResponse, error) {
	command := NewGetMatchmakingUserCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetMatchmakingUsers(ctx context.Context, in *api.GetMatchmakingUsersRequest) (*api.GetMatchmakingUsersResponse, error) {
	command := NewGetMatchmakingUsersCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateMatchmakingUser(ctx context.Context, in *api.UpdateMatchmakingUserRequest) (*api.UpdateMatchmakingUserResponse, error) {
	command := NewUpdateMatchmakingUserCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) SetMatchmakingUserElo(ctx context.Context, in *api.SetMatchmakingUserEloRequest) (*api.SetMatchmakingUserEloResponse, error) {
	command := NewSetMatchmakingUserEloCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) CreateMatchmakingTicket(ctx context.Context, in *api.CreateMatchmakingTicketRequest) (*api.CreateMatchmakingTicketResponse, error) {
	command := NewCreateMatchmakingTicketCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) PollMatchmakingTicket(ctx context.Context, in *api.MatchmakingTicketRequest) (*api.MatchmakingTicketResponse, error) {
	command := NewPollMatchmakingTicketCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetMatchmakingTicket(ctx context.Context, in *api.MatchmakingTicketRequest) (*api.GetMatchmakingTicketResponse, error) {
	command := NewGetMatchmakingTicketCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetMatchmakingTickets(ctx context.Context, in *api.GetMatchmakingTicketsRequest) (*api.GetMatchmakingTicketsResponse, error) {
	command := NewGetMatchmakingTicketsCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateMatchmakingTicket(ctx context.Context, in *api.UpdateMatchmakingTicketRequest) (*api.UpdateMatchmakingTicketResponse, error) {
	command := NewUpdateMatchmakingTicketCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) ExpireMatchmakingTicket(ctx context.Context, in *api.MatchmakingTicketRequest) (*api.ExpireMatchmakingTicketResponse, error) {
	command := NewExpireMatchmakingTicketCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) StartMatch(ctx context.Context, in *api.StartMatchRequest) (*api.StartMatchResponse, error) {
	command := NewStartMatchCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) EndMatch(ctx context.Context, in *api.EndMatchRequest) (*api.EndMatchResponse, error) {
	command := NewEndMatchCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetMatch(ctx context.Context, in *api.MatchRequest) (*api.GetMatchResponse, error) {
	command := NewGetMatchCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetMatches(ctx context.Context, in *api.GetMatchesRequest) (*api.GetMatchesResponse, error) {
	command := NewGetMatchesCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics).SetInvoker(invokers.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateMatch(ctx context.Context, in *api.UpdateMatchRequest) (*api.UpdateMatchResponse, error) {
	command := NewUpdateMatchCommand(s, in)
	invoker := invokers.NewLogInvoker().SetInvoker(invokers.NewTransportInvoker().SetInvoker(invokers.NewMetricsInvoker(s.metrics)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

/*

crud arena

crud user

set user elo
	id or user_id
	arena_id or name
	elo
	incrementElo bool

join queue
	- users array
		- id or user_id
	- arenas array
		- arena_id or name
	data

code:
	transaction
		atomic function that creates ticket with user_id and data if and only if there is no ticket with user_id that has
			- not expired
			- not in a match that is yet to finish
			- not in a match that is yet to start
		create ticket_arena rows with ticket_id and arena_id to satisfy many to many relationship
		create ticket_user rows with ticket_id and user_id to satisfy many to many relationship
		commit
	return ticket_id

poll queue
	- ticket_id

code:
	update ticket if ticket exists and is not expired and not in match
		- add extra time to the expiration time


leave queue
	- ticket_id

code:
	change expiration time to current time



update ticket
	- ticket_id
	- data

code:
	update ticket if ticket exists and is not expired and not in match
		- update data

check ticket
	- ticket_id

code:
	check if ticket has match_id
		- if yes, return
			- ticket_id
			- match
				- id
				- data
				- started_at
				- ended_at
				- created_at
				- updated_at
			- users
				- id
				- user_id
				- data
				- created_at
				- updated_at
			- data
			- created_at
			- updated_at
			- expires_at
		- if no, return ticket

start match
	- ticket_id or match_id
	- started_at

code:
	update match started at if match exists and is not started yet

end match
	- ticket_id or match_id
	- ended_at

code:
	update match ended at if match exists and is not ended yet

update match
	- ticket_id or match_id
	- data

code:
	update match data if match exists and is not ended yet

*/
