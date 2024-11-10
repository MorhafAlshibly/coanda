package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetMatchesCommand struct {
	service *Service
	In      *api.GetMatchesRequest
	Out     *api.GetMatchesResponse
}

func NewGetMatchesCommand(service *Service, in *api.GetMatchesRequest) *GetMatchesCommand {
	return &GetMatchesCommand{
		service: service,
		In:      in,
	}
}

func (c *GetMatchesCommand) Execute(ctx context.Context) error {
	// Check if matchmaking user is nil
	if c.In.MatchmakingUser == nil {
		c.In.MatchmakingUser = &api.MatchmakingUserRequest{}
	}
	// Check if arena is nil
	if c.In.Arena == nil {
		c.In.Arena = &api.ArenaRequest{}
	}
	matchLimit, matchOffset := conversion.PaginationToLimitOffset(c.In.Pagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	ticketLimit, ticketOffset := conversion.PaginationToLimitOffset(c.In.TicketPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	userLimit, userOffset := conversion.PaginationToLimitOffset(c.In.UserPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	// Matches hold a group of tickets, and tickets hold a group of users so overall limit is match limit * ticket limit * user limit
	limit := matchLimit * ticketLimit * userLimit
	// Matches hold a group of tickets, and tickets hold a group of users so overall offset is match offset * ticket limit * user limit + ticket offset * user limit + user offset
	offset := matchOffset*ticketLimit*userLimit + ticketOffset*userLimit + userOffset
	matches, err := c.service.database.GetMatches(ctx, model.GetMatchesParams{
		Arena: model.GetArenaParams{
			ID:   conversion.Uint64ToSqlNullInt64(c.In.Arena.Id),
			Name: conversion.StringToSqlNullString(c.In.Arena.Name),
		},
		MatchmakingUser: model.MatchmakingUserParams{
			ID:           conversion.Uint64ToSqlNullInt64(c.In.MatchmakingUser.Id),
			ClientUserID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingUser.ClientUserId),
		},
		Status: conversion.StringToSqlNullString(conversion.ValueToPointer(c.In.Status.String())),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return err
	}
	apiMatches, err := unmarshalMatches(matches)
	if err != nil {
		return err
	}
	c.Out = &api.GetMatchesResponse{
		Success: true,
		Matches: apiMatches,
	}
	return nil
}
