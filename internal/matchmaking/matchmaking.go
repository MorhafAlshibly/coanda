package matchmaking

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/cache"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/invoker"
	"github.com/MorhafAlshibly/coanda/pkg/metric"
)

type Service struct {
	api.UnimplementedMatchmakingServiceServer
	sql                  *sql.DB
	database             *model.Queries
	cache                cache.Cacher
	metric               metric.Metric
	minArenaNameLength   uint8
	maxArenaNameLength   uint8
	expiryTimeWindow     time.Duration
	startTimeBuffer      time.Duration
	lockedAtBuffer       time.Duration
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

func WithStartTimeBuffer(startTimeBuffer time.Duration) func(*Service) {
	return func(input *Service) {
		input.startTimeBuffer = startTimeBuffer
	}
}

func WithLockedAtBuffer(lockedAtBuffer time.Duration) func(*Service) {
	return func(input *Service) {
		input.lockedAtBuffer = lockedAtBuffer
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
		startTimeBuffer:      10 * time.Second,
		lockedAtBuffer:       5 * time.Second,
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
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetArena(ctx context.Context, in *api.ArenaRequest) (*api.GetArenaResponse, error) {
	command := NewGetArenaCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric).SetInvoker(invoker.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetArenas(ctx context.Context, in *api.Pagination) (*api.GetArenasResponse, error) {
	command := NewGetArenasCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric).SetInvoker(invoker.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateArena(ctx context.Context, in *api.UpdateArenaRequest) (*api.UpdateArenaResponse, error) {
	command := NewUpdateArenaCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) CreateMatchmakingUser(ctx context.Context, in *api.CreateMatchmakingUserRequest) (*api.CreateMatchmakingUserResponse, error) {
	command := NewCreateMatchmakingUserCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetMatchmakingUser(ctx context.Context, in *api.MatchmakingUserRequest) (*api.GetMatchmakingUserResponse, error) {
	command := NewGetMatchmakingUserCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric).SetInvoker(invoker.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetMatchmakingUsers(ctx context.Context, in *api.Pagination) (*api.GetMatchmakingUsersResponse, error) {
	command := NewGetMatchmakingUsersCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric).SetInvoker(invoker.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateMatchmakingUser(ctx context.Context, in *api.UpdateMatchmakingUserRequest) (*api.UpdateMatchmakingUserResponse, error) {
	command := NewUpdateMatchmakingUserCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) CreateMatchmakingTicket(ctx context.Context, in *api.CreateMatchmakingTicketRequest) (*api.CreateMatchmakingTicketResponse, error) {
	command := NewCreateMatchmakingTicketCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) PollMatchmakingTicket(ctx context.Context, in *api.GetMatchmakingTicketRequest) (*api.GetMatchmakingTicketResponse, error) {
	command := NewPollMatchmakingTicketCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetMatchmakingTicket(ctx context.Context, in *api.GetMatchmakingTicketRequest) (*api.GetMatchmakingTicketResponse, error) {
	command := NewGetMatchmakingTicketCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric).SetInvoker(invoker.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetMatchmakingTickets(ctx context.Context, in *api.GetMatchmakingTicketsRequest) (*api.GetMatchmakingTicketsResponse, error) {
	command := NewGetMatchmakingTicketsCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric).SetInvoker(invoker.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateMatchmakingTicket(ctx context.Context, in *api.UpdateMatchmakingTicketRequest) (*api.UpdateMatchmakingTicketResponse, error) {
	command := NewUpdateMatchmakingTicketCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) ExpireMatchmakingTicket(ctx context.Context, in *api.MatchmakingTicketRequest) (*api.ExpireMatchmakingTicketResponse, error) {
	command := NewExpireMatchmakingTicketCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) StartMatch(ctx context.Context, in *api.StartMatchRequest) (*api.StartMatchResponse, error) {
	command := NewStartMatchCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) EndMatch(ctx context.Context, in *api.EndMatchRequest) (*api.EndMatchResponse, error) {
	command := NewEndMatchCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetMatch(ctx context.Context, in *api.GetMatchRequest) (*api.GetMatchResponse, error) {
	command := NewGetMatchCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric).SetInvoker(invoker.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) GetMatches(ctx context.Context, in *api.GetMatchesRequest) (*api.GetMatchesResponse, error) {
	command := NewGetMatchesCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric).SetInvoker(invoker.NewCacheInvoker(s.cache))))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) UpdateMatch(ctx context.Context, in *api.UpdateMatchRequest) (*api.UpdateMatchResponse, error) {
	command := NewUpdateMatchCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
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

func unmarshalMatchmakingUser(matchmakingUser model.MatchmakingUser) (*api.MatchmakingUser, error) {
	data, err := conversion.RawJsonToProtobufStruct(matchmakingUser.Data)
	if err != nil {
		return nil, err
	}
	return &api.MatchmakingUser{
		Id:           matchmakingUser.ID,
		ClientUserId: matchmakingUser.ClientUserID,
		Elo:          matchmakingUser.Elo,
		Data:         data,
		CreatedAt:    conversion.TimeToTimestamppb(&matchmakingUser.CreatedAt),
		UpdatedAt:    conversion.TimeToTimestamppb(&matchmakingUser.UpdatedAt),
	}, nil
}

func unmarshalMatchmakingTicket(matchmakingTicket []model.MatchmakingTicketWithUserAndArena) (*api.MatchmakingTicket, error) {
	data, err := conversion.RawJsonToProtobufStruct(matchmakingTicket[0].TicketData)
	if err != nil {
		return nil, err
	}
	out := &api.MatchmakingTicket{
		Id:        matchmakingTicket[0].TicketID,
		MatchId:   conversion.SqlNullInt64ToUint64(matchmakingTicket[0].MatchmakingMatchID),
		Status:    api.MatchmakingTicket_Status(api.MatchmakingTicket_Status_value[matchmakingTicket[0].Status]),
		Data:      data,
		ExpiresAt: conversion.TimeToTimestamppb(&matchmakingTicket[0].ExpiresAt),
		CreatedAt: conversion.TimeToTimestamppb(&matchmakingTicket[0].TicketCreatedAt),
		UpdatedAt: conversion.TimeToTimestamppb(&matchmakingTicket[0].TicketUpdatedAt),
	}
	users := []*api.MatchmakingUser{}
	arenas := []*api.Arena{}
	currentUserId := uint64(0)
	currentArenaId := uint64(0)
	for _, ticket := range matchmakingTicket {
		if ticket.MatchmakingUserID != currentUserId {
			currentUserId = ticket.MatchmakingUserID
			user, err := unmarshalMatchmakingUser(model.MatchmakingUser{
				ID:           ticket.MatchmakingUserID,
				ClientUserID: ticket.ClientUserID,
				Elo:          ticket.Elo,
				Data:         ticket.UserData,
				CreatedAt:    ticket.UserCreatedAt,
				UpdatedAt:    ticket.UserUpdatedAt,
			})
			if err != nil {
				return nil, err
			}
			users = append(users, user)
		}
		if ticket.ArenaID != currentArenaId {
			currentArenaId = ticket.ArenaID
			arena, err := unmarshalArena(model.MatchmakingArena{
				ID:                  ticket.ArenaID,
				Name:                ticket.ArenaName,
				MinPlayers:          ticket.ArenaMinPlayers,
				MaxPlayersPerTicket: ticket.ArenaMaxPlayersPerTicket,
				MaxPlayers:          ticket.ArenaMaxPlayers,
				Data:                ticket.ArenaData,
				CreatedAt:           ticket.ArenaCreatedAt,
				UpdatedAt:           ticket.ArenaUpdatedAt,
			})
			if err != nil {
				return nil, err
			}
			arenas = append(arenas, arena)
		}
	}
	out.MatchmakingUsers = users
	out.Arenas = arenas
	return out, nil
}

func unmarshalMatchmakingTickets(matchmakingTickets []model.MatchmakingTicketWithUserAndArena) ([]*api.MatchmakingTicket, error) {
	// Tickets are already sorted by ticket ID
	unmarshalledTickets := make([]*api.MatchmakingTicket, 0)
	for i := 0; i < len(matchmakingTickets); {
		ticket := make([]model.MatchmakingTicketWithUserAndArena, 0)
		for j := i; j < len(matchmakingTickets) && matchmakingTickets[j].TicketID == matchmakingTickets[i].TicketID; j++ {
			ticket = append(ticket, matchmakingTickets[j])
			i++
		}
		fmt.Printf("ticket: %+v\n", ticket)
		unmarshalledTicket, err := unmarshalMatchmakingTicket(ticket)
		if err != nil {
			return nil, err
		}
		unmarshalledTickets = append(unmarshalledTickets, unmarshalledTicket)
	}
	return unmarshalledTickets, nil
}

func unmarshalMatch(match []model.MatchmakingMatchWithArenaAndTicket) (*api.Match, error) {
	data, err := conversion.RawJsonToProtobufStruct(match[0].MatchData)
	if err != nil {
		return nil, err
	}
	arena, err := unmarshalArena(model.MatchmakingArena{
		ID:                  uint64(match[0].ArenaID.Int64),
		Name:                match[0].ArenaName.String,
		MinPlayers:          uint32(match[0].ArenaMinPlayers.Int32),
		MaxPlayersPerTicket: uint32(match[0].ArenaMaxPlayersPerTicket.Int32),
		MaxPlayers:          uint32(match[0].ArenaMaxPlayers.Int32),
		Data:                match[0].ArenaData,
		CreatedAt:           match[0].ArenaCreatedAt.Time,
		UpdatedAt:           match[0].ArenaUpdatedAt.Time,
	})
	if err != nil {
		return nil, err
	}
	out := &api.Match{
		Id:              uint64(match[0].MatchID.Int64),
		Arena:           arena,
		PrivateServerId: conversion.SqlNullStringToString(match[0].PrivateServerID),
		Status:          api.Match_Status(api.Match_Status_value[match[0].MatchStatus.String]),
		Data:            data,
		LockedAt:        conversion.TimeToTimestamppb(&match[0].LockedAt.Time),
		StartedAt:       conversion.TimeToTimestamppb(&match[0].StartedAt.Time),
		EndedAt:         conversion.TimeToTimestamppb(&match[0].EndedAt.Time),
		CreatedAt:       conversion.TimeToTimestamppb(&match[0].MatchCreatedAt.Time),
		UpdatedAt:       conversion.TimeToTimestamppb(&match[0].MatchUpdatedAt.Time),
	}
	tickets := []model.MatchmakingTicketWithUserAndArena{}
	for _, ticket := range match {
		tickets = append(tickets, model.MatchmakingTicketWithUserAndArena{
			TicketID:                 uint64(ticket.TicketID.Int64),
			MatchmakingMatchID:       ticket.MatchID,
			Status:                   ticket.TicketStatus.String,
			TicketData:               ticket.TicketData,
			ExpiresAt:                ticket.ExpiresAt.Time,
			TicketCreatedAt:          ticket.TicketCreatedAt.Time,
			TicketUpdatedAt:          ticket.TicketUpdatedAt.Time,
			MatchmakingUserID:        uint64(ticket.MatchmakingUserID.Int64),
			ClientUserID:             uint64(ticket.ClientUserID.Int64),
			Elo:                      ticket.Elo.Int64,
			UserNumber:               ticket.UserNumber,
			UserData:                 ticket.UserData,
			UserCreatedAt:            ticket.UserCreatedAt.Time,
			UserUpdatedAt:            ticket.UserUpdatedAt.Time,
			ArenaID:                  uint64(ticket.TicketArenaID.Int64),
			ArenaName:                ticket.TicketArenaName.String,
			ArenaMinPlayers:          uint32(ticket.TicketArenaMinPlayers.Int32),
			ArenaMaxPlayersPerTicket: uint32(ticket.TicketArenaMaxPlayersPerTicket.Int32),
			ArenaMaxPlayers:          uint32(ticket.TicketArenaMaxPlayers.Int32),
			ArenaNumber:              ticket.ArenaNumber,
			ArenaData:                ticket.TicketArenaData,
			ArenaCreatedAt:           ticket.TicketArenaCreatedAt.Time,
			ArenaUpdatedAt:           ticket.TicketArenaUpdatedAt.Time,
		})
	}
	unmarshalledTickets, err := unmarshalMatchmakingTickets(tickets)
	if err != nil {
		return nil, err
	}
	out.Tickets = unmarshalledTickets
	return out, nil
}

func unmarshalMatches(matches []model.MatchmakingMatchWithArenaAndTicket) ([]*api.Match, error) {
	// Matches are already sorted by match ID
	unmarshalledMatches := make([]*api.Match, 0)
	for i := 0; i < len(matches); {
		match := make([]model.MatchmakingMatchWithArenaAndTicket, 0)
		for j := i; j < len(matches) && matches[j].MatchID == matches[i].MatchID; j++ {
			match = append(match, matches[j])
			i++
		}
		unmarshalledMatch, err := unmarshalMatch(match)
		if err != nil {
			return nil, err
		}
		unmarshalledMatches = append(unmarshalledMatches, unmarshalledMatch)
	}
	return unmarshalledMatches, nil
}

// Enum for errors
type MatchmakingRequestError string

const (
	NAME_TOO_SHORT                                 MatchmakingRequestError = "NAME_TOO_SHORT"
	NAME_TOO_LONG                                  MatchmakingRequestError = "NAME_TOO_LONG"
	ID_OR_NAME_REQUIRED                            MatchmakingRequestError = "ID_OR_NAME_REQUIRED"
	MATCHMAKING_USER_ID_OR_CLIENT_USER_ID_REQUIRED MatchmakingRequestError = "MATCHMAKING_USER_ID_OR_USER_ID_REQUIRED"
	TICKET_ID_OR_MATCHMAKING_USER_REQUIRED         MatchmakingRequestError = "TICKET_ID_OR_MATCHMAKING_USER_REQUIRED"
	ID_OR_MATCHMAKING_TICKET_REQUIRED              MatchmakingRequestError = "ID_OR_MATCHMAKING_TICKET_REQUIRED"
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
		return conversion.ValueToPointer(TICKET_ID_OR_MATCHMAKING_USER_REQUIRED)
	}
	if request.Id != nil {
		return nil
	}
	if request.MatchmakingUser == nil {
		return conversion.ValueToPointer(TICKET_ID_OR_MATCHMAKING_USER_REQUIRED)
	}
	return s.checkForMatchmakingUserRequestError(request.MatchmakingUser)
}

func (s *Service) checkForMatchRequestError(request *api.MatchRequest) *MatchmakingRequestError {
	if request == nil {
		return conversion.ValueToPointer(ID_OR_MATCHMAKING_TICKET_REQUIRED)
	}
	if request.Id != nil {
		return nil
	}
	if request.MatchmakingTicket == nil {
		return conversion.ValueToPointer(ID_OR_MATCHMAKING_TICKET_REQUIRED)
	}
	return s.checkForMatchmakingTicketRequestError(request.MatchmakingTicket)
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
