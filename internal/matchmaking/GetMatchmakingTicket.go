package matchmaking

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetMatchmakingTicketCommand struct {
	service *Service
	In      *api.GetMatchmakingTicketRequest
	Out     *api.GetMatchmakingTicketResponse
}

func NewGetMatchmakingTicketCommand(service *Service, in *api.GetMatchmakingTicketRequest) *GetMatchmakingTicketCommand {
	return &GetMatchmakingTicketCommand{
		service: service,
		In:      in,
	}
}

func (c *GetMatchmakingTicketCommand) Execute(ctx context.Context) error {
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
	limit, offset := conversion.PaginationToLimitOffset(c.In.Pagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	matchmakingTicket, err := c.service.database.GetMatchmakingTicket(ctx, model.GetMatchmakingTicketParams{
		MatchmakingTicket: model.MatchmakingTicketParams{
			MatchmakingUser: model.MatchmakingUserParams{
				ID:           conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.Id),
				ClientUserID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.MatchmakingUser.ClientUserId),
			},
			ID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.Id),
		},
		Limit:  limit,
		Offset: offset,
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
