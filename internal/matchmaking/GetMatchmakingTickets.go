package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetMatchmakingTicketsCommand struct {
	service *Service
	In      *api.GetMatchmakingTicketsRequest
	Out     *api.GetMatchmakingTicketsResponse
}

func NewGetMatchmakingTicketsCommand(service *Service, in *api.GetMatchmakingTicketsRequest) *GetMatchmakingTicketsCommand {
	return &GetMatchmakingTicketsCommand{
		service: service,
		In:      in,
	}
}

func (c *GetMatchmakingTicketsCommand) Execute(ctx context.Context) error {
	limit, offset := conversion.PaginationToLimitOffset(c.In.Pagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	userLimit, userOffset := conversion.PaginationToLimitOffset(c.In.UserPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	arenaLimit, arenaOffset := conversion.PaginationToLimitOffset(c.In.ArenaPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	statuses := []string{}
	for _, status := range c.In.Statuses {
		statuses = append(statuses, status.String())
	}
	tickets, err := c.service.database.GetMatchmakingTickets(ctx, model.GetMatchmakingTicketsParams{
		MatchmakingMatchID: conversion.Uint64ToSqlNullInt64(c.In.MatchId),
		MatchmakingUser:    matchmakingUserRequestToMatchmakingUserParams(c.In.MatchmakingUser),
		Statuses:           statuses,
		Limit:              limit,
		Offset:             offset,
		UserLimit:          userLimit,
		UserOffset:         userOffset,
		ArenaLimit:         arenaLimit,
		ArenaOffset:        arenaOffset,
	})
	if err != nil {
		return err
	}
	apiTickets, err := unmarshalMatchmakingTickets(tickets)
	if err != nil {
		return err
	}
	c.Out = &api.GetMatchmakingTicketsResponse{
		Success:            true,
		MatchmakingTickets: apiTickets,
	}
	return nil
}
