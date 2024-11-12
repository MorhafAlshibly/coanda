package matchmaking

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type PollMatchmakingTicketCommand struct {
	service *Service
	In      *api.GetMatchmakingTicketRequest
	Out     *api.GetMatchmakingTicketResponse
}

func NewPollMatchmakingTicketCommand(service *Service, in *api.GetMatchmakingTicketRequest) *PollMatchmakingTicketCommand {
	return &PollMatchmakingTicketCommand{
		service: service,
		In:      in,
	}
}

func (c *PollMatchmakingTicketCommand) Execute(ctx context.Context) error {
	mtErr := c.service.checkForMatchmakingTicketRequestError(c.In.MatchmakingTicket)
	if mtErr != nil {
		c.Out = &api.GetMatchmakingTicketResponse{
			Success: false,
			Error:   conversion.Enum(*mtErr, api.GetMatchmakingTicketResponse_Error_value, api.GetMatchmakingTicketResponse_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED),
		}
		return nil
	}
	// Make sure matchmaking user isnt nil
	if c.In.MatchmakingTicket.MatchmakingUser == nil {
		c.In.MatchmakingTicket.MatchmakingUser = &api.MatchmakingUserRequest{}
	}
	// Start transaction
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	_, err = qtx.PollMatchmakingTicket(ctx, model.PollMatchmakingTicketParams{
		MatchmakingTicket: model.MatchmakingTicketParams{
			MatchmakingUser: model.MatchmakingUserParams{
				ID:           conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.Id),
				ClientUserID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.MatchmakingUser.ClientUserId),
			},
			ID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.Id),
			// Only update tickets that are pending
			Statuses: []string{"PENDING"},
		},
		ExpiryTimeWindow: c.service.expiryTimeWindow,
	})
	if err != nil {
		return err
	}
	userLimit, userOffset := conversion.PaginationToLimitOffset(c.In.UserPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	arenaLimit, arenaOffset := conversion.PaginationToLimitOffset(c.In.ArenaPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	matchmakingTicket, err := qtx.GetMatchmakingTicket(ctx, model.GetMatchmakingTicketParams{
		MatchmakingTicket: model.MatchmakingTicketParams{
			MatchmakingUser: model.MatchmakingUserParams{
				ID:           conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.Id),
				ClientUserID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.MatchmakingUser.ClientUserId),
			},
			ID:       conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.Id),
			Statuses: []string{"PENDING", "MATCHED"},
		},
		UserLimit:   userLimit,
		UserOffset:  userOffset,
		ArenaLimit:  arenaLimit,
		ArenaOffset: arenaOffset,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			c.Out = &api.GetMatchmakingTicketResponse{
				Success: false,
				Error:   api.GetMatchmakingTicketResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	apiMatchmakingTicket, err := unmarshalMatchmakingTicket(matchmakingTicket)
	if err != nil {
		return err
	}
	c.Out = &api.GetMatchmakingTicketResponse{
		Success:           true,
		MatchmakingTicket: apiMatchmakingTicket,
		Error:             api.GetMatchmakingTicketResponse_NONE,
	}
	return nil
}
