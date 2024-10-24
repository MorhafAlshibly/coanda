package matchmaking

import (
	"context"

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
	result, err := qtx.PollMatchmakingTicket(ctx, model.PollMatchmakingTicketParams{
		MatchmakingTicket: model.MatchmakingTicketParams{
			MatchmakingUser: model.GetMatchmakingUserParams{
				ID:           conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.Id),
				ClientUserID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.MatchmakingUser.ClientUserId),
			},
			ID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.Id),
		},
		ExpiryTimeWindow: c.service.expiryTimeWindow,
	})
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// We didn't find a row, cant be an unchanged row because we are updating the expiry time based on the current time
		c.Out = &api.GetMatchmakingTicketResponse{
			Success: false,
			Error:   api.GetMatchmakingTicketResponse_NOT_FOUND,
		}
		return nil
	}
	limit, offset := conversion.PaginationToLimitOffset(c.In.Pagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	matchmakingTicket, err := qtx.GetMatchmakingTicket(ctx, model.GetMatchmakingTicketParams{
		MatchmakingTicket: model.MatchmakingTicketParams{
			MatchmakingUser: model.GetMatchmakingUserParams{
				ID:           conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.Id),
				ClientUserID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.MatchmakingUser.ClientUserId),
			},
			ID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.Id),
		},
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
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
