package matchmaking

import (
	"context"
	"database/sql"
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

func (s *Service) DeleteMatchmakingTicket(ctx context.Context, in *api.MatchmakingTicketRequest) (*api.DeleteMatchmakingTicketResponse, error) {
	command := NewDeleteMatchmakingTicketCommand(s, in)
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

func (s *Service) SetMatchPrivateServer(ctx context.Context, in *api.SetMatchPrivateServerRequest) (*api.SetMatchPrivateServerResponse, error) {
	command := NewSetMatchPrivateServerCommand(s, in)
	invoker := invoker.NewLogInvoker().SetInvoker(invoker.NewTransportInvoker().SetInvoker(invoker.NewMetricInvoker(s.metric)))
	err := invoker.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}
	return command.Out, nil
}

func (s *Service) DeleteMatch(ctx context.Context, in *api.MatchRequest) (*api.DeleteMatchResponse, error) {
	command := NewDeleteMatchCommand(s, in)
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
		CreatedAt: conversion.TimeToTimestamppb(&matchmakingTicket[0].TicketCreatedAt),
		UpdatedAt: conversion.TimeToTimestamppb(&matchmakingTicket[0].TicketUpdatedAt),
	}
	users := []*api.MatchmakingUser{}
	arenas := []*api.Arena{}
	currentUserId := uint64(0)
	// To avoid duplicate arenas, because we are sorting by ticket, then user, so we only need to check if we have seen the arena ID before
	seenArenaIds := map[uint64]bool{}
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
		if _, ok := seenArenaIds[ticket.ArenaID]; !ok {
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
			seenArenaIds[ticket.ArenaID] = true
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
		lastTicketID := matchmakingTickets[i].TicketID
		for j := i; j < len(matchmakingTickets) && matchmakingTickets[j].TicketID == lastTicketID; j++ {
			ticket = append(ticket, matchmakingTickets[j])
			i++
		}
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
		ID:                  uint64(match[0].ArenaID),
		Name:                match[0].ArenaName,
		MinPlayers:          uint32(match[0].ArenaMinPlayers),
		MaxPlayersPerTicket: uint32(match[0].ArenaMaxPlayersPerTicket),
		MaxPlayers:          uint32(match[0].ArenaMaxPlayers),
		Data:                match[0].ArenaData,
		CreatedAt:           match[0].ArenaCreatedAt,
		UpdatedAt:           match[0].ArenaUpdatedAt,
	})
	if err != nil {
		return nil, err
	}
	out := &api.Match{
		Id:              uint64(match[0].MatchID),
		Arena:           arena,
		PrivateServerId: conversion.SqlNullStringToString(match[0].PrivateServerID),
		Status:          api.Match_Status(api.Match_Status_value[match[0].MatchStatus]),
		Data:            data,
		LockedAt:        conversion.TimeToTimestamppb(&match[0].LockedAt.Time),
		StartedAt:       conversion.TimeToTimestamppb(&match[0].StartedAt.Time),
		EndedAt:         conversion.TimeToTimestamppb(&match[0].EndedAt.Time),
		CreatedAt:       conversion.TimeToTimestamppb(&match[0].MatchCreatedAt),
		UpdatedAt:       conversion.TimeToTimestamppb(&match[0].MatchUpdatedAt),
	}
	tickets := []model.MatchmakingTicketWithUserAndArena{}
	for _, ticket := range match {
		tickets = append(tickets, model.MatchmakingTicketWithUserAndArena{
			TicketID:                 uint64(ticket.TicketID.Int64),
			MatchmakingMatchID:       conversion.Uint64ToSqlNullInt64(&ticket.MatchID),
			Status:                   ticket.TicketStatus.String,
			TicketData:               ticket.TicketData,
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
		lastMatchID := matches[i].MatchID
		for j := i; j < len(matches) && matches[j].MatchID == lastMatchID; j++ {
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

func arenaRequestToArenaParams(request *api.ArenaRequest) model.ArenaParams {
	if request == nil {
		return model.ArenaParams{}
	}
	return model.ArenaParams{
		ID:   conversion.Uint64ToSqlNullInt64(request.Id),
		Name: conversion.StringToSqlNullString(request.Name),
	}
}

func matchmakingUserRequestToMatchmakingUserParams(request *api.MatchmakingUserRequest) model.MatchmakingUserParams {
	if request == nil {
		return model.MatchmakingUserParams{}
	}
	return model.MatchmakingUserParams{
		ID:           conversion.Uint64ToSqlNullInt64(request.Id),
		ClientUserID: conversion.Uint64ToSqlNullInt64(request.ClientUserId),
	}
}

func matchmakingTicketRequestToMatchmakingTicketParams(request *api.MatchmakingTicketRequest) model.MatchmakingTicketParams {
	if request == nil {
		return model.MatchmakingTicketParams{}
	}
	return model.MatchmakingTicketParams{
		ID:              conversion.Uint64ToSqlNullInt64(request.Id),
		MatchmakingUser: matchmakingUserRequestToMatchmakingUserParams(request.MatchmakingUser),
	}
}

func matchRequestToMatchParams(request *api.MatchRequest) model.MatchParams {
	if request == nil {
		return model.MatchParams{}
	}
	return model.MatchParams{
		ID:                conversion.Uint64ToSqlNullInt64(request.Id),
		MatchmakingTicket: matchmakingTicketRequestToMatchmakingTicketParams(request.MatchmakingTicket),
	}
}

// Enum for errors
type MatchmakingRequestError string

const (
	NAME_TOO_SHORT                                     MatchmakingRequestError = "NAME_TOO_SHORT"
	NAME_TOO_LONG                                      MatchmakingRequestError = "NAME_TOO_LONG"
	ARENA_ID_OR_NAME_REQUIRED                          MatchmakingRequestError = "ARENA_ID_OR_NAME_REQUIRED"
	MATCHMAKING_USER_ID_OR_CLIENT_USER_ID_REQUIRED     MatchmakingRequestError = "MATCHMAKING_USER_ID_OR_CLIENT_USER_ID_REQUIRED"
	MATCHMAKING_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED MatchmakingRequestError = "MATCHMAKING_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED"
	MATCH_ID_OR_MATCHMAKING_TICKET_REQUIRED            MatchmakingRequestError = "MATCH_ID_OR_MATCHMAKING_TICKET_REQUIRED"
)

func (s *Service) checkForArenaRequestError(request *api.ArenaRequest) *MatchmakingRequestError {
	if request == nil {
		return conversion.ValueToPointer(ARENA_ID_OR_NAME_REQUIRED)
	}
	if request.Id != nil {
		return nil
	}
	if request.Name == nil {
		return conversion.ValueToPointer(ARENA_ID_OR_NAME_REQUIRED)
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
		return conversion.ValueToPointer(MATCHMAKING_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED)
	}
	if request.Id != nil {
		return nil
	}
	if request.MatchmakingUser == nil {
		return conversion.ValueToPointer(MATCHMAKING_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED)
	}
	return s.checkForMatchmakingUserRequestError(request.MatchmakingUser)
}

func (s *Service) checkForMatchRequestError(request *api.MatchRequest) *MatchmakingRequestError {
	if request == nil {
		return conversion.ValueToPointer(MATCH_ID_OR_MATCHMAKING_TICKET_REQUIRED)
	}
	if request.Id != nil {
		return nil
	}
	if request.MatchmakingTicket == nil {
		return conversion.ValueToPointer(MATCH_ID_OR_MATCHMAKING_TICKET_REQUIRED)
	}
	return s.checkForMatchmakingTicketRequestError(request.MatchmakingTicket)
}
