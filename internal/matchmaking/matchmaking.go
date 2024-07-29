package matchmaking

import (
	"context"
	"database/sql"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
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
	expiryTimeWindow     time.Duration
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

func WithExpiryTimeWindow(expiryTimeWindow time.Duration) func(*Service) {
	return func(input *Service) {
		input.expiryTimeWindow = expiryTimeWindow
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
		minArenaNameLength:   3,
		maxArenaNameLength:   20,
		expiryTimeWindow:     5 * time.Second,
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

func (s *Service) GetArenas(ctx context.Context, in *api.Pagination) (*api.GetArenasResponse, error) {
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

func (s *Service) GetMatchmakingUsers(ctx context.Context, in *api.Pagination) (*api.GetMatchmakingUsersResponse, error) {
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

func (s *Service) GetMatchmakingTicket(ctx context.Context, in *api.GetMatchmakingTicketRequest) (*api.GetMatchmakingTicketResponse, error) {
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

func unmarshalArena(arena model.MatchmakingArena) (*api.Arena, error) {
	data, err := conversion.RawJsonToProtobufStruct(arena.Data)
	if err != nil {
		return nil, err
	}
	return &api.Arena{
		Id:                  arena.ID,
		Name:                arena.Name,
		MinPlayers:          uint32(arena.MinPlayers),
		MaxPlayersPerTicket: uint32(arena.MaxPlayersPerTicket),
		MaxPlayers:          uint32(arena.MaxPlayers),
		Data:                data,
		CreatedAt:           conversion.TimeToTimestamppb(&arena.CreatedAt),
		UpdatedAt:           conversion.TimeToTimestamppb(&arena.UpdatedAt),
	}, nil
}

func unmarshalMatchmakingUser(matchmakingUser model.MatchmakingUserWithElo) (*api.MatchmakingUser, error) {
	data, err := conversion.RawJsonToProtobufStruct(matchmakingUser.Data)
	if err != nil {
		return nil, err
	}
	// Convert elos to array of map with keys arena ID and elo
	elos, err := conversion.RawJsonToArrayOfMaps(matchmakingUser.Elos)
	if err != nil {
		return nil, err
	}
	eloObjects := make([]*api.MatchmakingUserElo, len(elos))
	for i, elo := range elos {
		eloObjects[i] = &api.MatchmakingUserElo{
			ArenaId: uint64(elo["arena_id"].(int64)),
			Elo:     int64(elo["elo"].(int64)),
		}
	}
	return &api.MatchmakingUser{
		Id:           matchmakingUser.ID,
		ClientUserId: matchmakingUser.ClientUserID,
		Data:         data,
		Elos:         eloObjects,
		CreatedAt:    conversion.TimeToTimestamppb(&matchmakingUser.CreatedAt),
		UpdatedAt:    conversion.TimeToTimestamppb(&matchmakingUser.UpdatedAt),
	}, nil
}

func unmarshalMatchmakingTicket(matchmakingTicket []model.MatchmakingTicketWithUserAndArena) (*api.MatchmakingTicket, error) {
	data, err := conversion.RawJsonToProtobufStruct(matchmakingTicket[0].TicketData)
	if err != nil {
		return nil, err
	}
	arenas, err := conversion.RawJsonToArrayOfMaps(matchmakingTicket[0].Arenas)
	if err != nil {
		return nil, err
	}
	arenaObjects := make([]*api.Arena, len(arenas))
	for i, arena := range arenas {
		arenaObjects[i] = &api.Arena{
			Id:                  uint64(arena["arena_id"].(int64)),
			Name:                arena["arena_name"].(string),
			MinPlayers:          uint32(arena["arena_min_players"].(int64)),
			MaxPlayersPerTicket: uint32(arena["arena_max_players_per_ticket"].(int64)),
			MaxPlayers:          uint32(arena["arena_max_players"].(int64)),
			CreatedAt:           conversion.TimeToTimestamppb(arena["arena_created_at"].(*time.Time)),
			UpdatedAt:           conversion.TimeToTimestamppb(arena["arena_updated_at"].(*time.Time)),
		}
	}
	users := make([]*api.MatchmakingUser, len(matchmakingTicket))
	for i, user := range matchmakingTicket {
		matchmakingUserWithElo := model.MatchmakingUserWithElo{
			ID:           user.MatchmakingUserID,
			ClientUserID: user.ClientUserID,
			Elos:         user.Elos,
			Data:         user.UserData,
			CreatedAt:    user.UserCreatedAt,
			UpdatedAt:    user.UserUpdatedAt,
		}
		matchmakingUser, err := unmarshalMatchmakingUser(matchmakingUserWithElo)
		if err != nil {
			return nil, err
		}
		users[i] = matchmakingUser
	}
	return &api.MatchmakingTicket{
		Id:               matchmakingTicket[0].ID,
		MatchmakingUsers: users,
		Arenas:           arenaObjects,
		MatchId:          conversion.SqlNullInt64ToUint64(matchmakingTicket[0].MatchmakingMatchID),
		Status:           api.MatchmakingTicket_Status(api.MatchmakingTicket_Status_value[matchmakingTicket[0].Status]),
		Data:             data,
		ExpiresAt:        conversion.TimeToTimestamppb(&matchmakingTicket[0].ExpiresAt),
		CreatedAt:        conversion.TimeToTimestamppb(&matchmakingTicket[0].TicketCreatedAt),
		UpdatedAt:        conversion.TimeToTimestamppb(&matchmakingTicket[0].TicketUpdatedAt),
	}, nil
}

// Enum for errors
type MatchmakingRequestError string

const (
	NAME_TOO_SHORT                                 MatchmakingRequestError = "NAME_TOO_SHORT"
	NAME_TOO_LONG                                  MatchmakingRequestError = "NAME_TOO_LONG"
	ID_OR_NAME_REQUIRED                            MatchmakingRequestError = "ID_OR_NAME_REQUIRED"
	MATCHMAKING_USER_ID_OR_CLIENT_USER_ID_REQUIRED MatchmakingRequestError = "MATCHMAKING_USER_ID_OR_USER_ID_REQUIRED"
	ID_OR_MATCHMAKING_USER_REQUIRED                MatchmakingRequestError = "ID_OR_MATCHMAKING_USER_REQUIRED"
)

func (s *Service) checkForArenaRequestError(request *api.ArenaRequest) *MatchmakingRequestError {
	if request == nil {
		return conversion.ValueToPointer(ID_OR_NAME_REQUIRED)
	}
	if request.Id != nil {
		return nil
	}
	if request.Name == nil {
		return conversion.ValueToPointer(ID_OR_NAME_REQUIRED)
	}
	if len(*request.Name) < int(s.minArenaNameLength) {
		return conversion.ValueToPointer(NAME_TOO_SHORT)
	}
	if len(*request.Name) > int(s.maxArenaNameLength) {
		return conversion.ValueToPointer(NAME_TOO_LONG)
	}
	return nil
}

func (s *Service) checkForMatchmakingUserRequestError(request *api.MatchmakingUserRequest) *MatchmakingRequestError {
	if request == nil {
		return conversion.ValueToPointer(MATCHMAKING_USER_ID_OR_CLIENT_USER_ID_REQUIRED)
	}
	if request.Id != nil {
		return nil
	}
	if request.ClientUserId == nil {
		return conversion.ValueToPointer(MATCHMAKING_USER_ID_OR_CLIENT_USER_ID_REQUIRED)
	}
	return nil
}

func (s *Service) checkForMatchmakingTicketRequestError(request *api.MatchmakingTicketRequest) *MatchmakingRequestError {
	if request == nil {
		return conversion.ValueToPointer(ID_OR_MATCHMAKING_USER_REQUIRED)
	}
	if request.Id != nil {
		return nil
	}
	if request.MatchmakingUser == nil {
		return conversion.ValueToPointer(ID_OR_MATCHMAKING_USER_REQUIRED)
	}
	return s.checkForMatchmakingUserRequestError(request.MatchmakingUser)
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
