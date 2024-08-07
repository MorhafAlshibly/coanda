package matchmaking

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetMatchCommand struct {
	service *Service
	In      *api.GetMatchRequest
	Out     *api.GetMatchResponse
}

func NewGetMatchCommand(service *Service, in *api.GetMatchRequest) *GetMatchCommand {
	return &GetMatchCommand{
		service: service,
		In:      in,
	}
}

func (c *GetMatchCommand) Execute(ctx context.Context) error {
	mmErr := c.service.checkForMatchRequestError(c.In.Match)
	if mmErr != nil {
		c.Out = &api.GetMatchResponse{
			Success: false,
			Error:   conversion.Enum(*mmErr, api.GetMatchResponse_Error_value, api.GetMatchResponse_ID_OR_MATCHMAKING_TICKET_REQUIRED),
		}
		return nil
	}
	// Make sure matchmaking ticket isnt nil
	if c.In.Match.MatchmakingTicket == nil {
		c.In.Match.MatchmakingTicket = &api.MatchmakingTicketRequest{
			MatchmakingUser: &api.MatchmakingUserRequest{},
		}
	}
	// Make sure matchmaking user isnt nil
	if c.In.Match.MatchmakingTicket.MatchmakingUser == nil {
		c.In.Match.MatchmakingTicket.MatchmakingUser = &api.MatchmakingUserRequest{}
	}
	ticketLimit, ticketOffset := conversion.PaginationToLimitOffset(c.In.TicketPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	userLimit, userOffset := conversion.PaginationToLimitOffset(c.In.UserPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	// Tickets hold a limit of users so overall limit is ticket limit * user limit
	limit := ticketLimit * userLimit
	// Tickets hold a limit of users so overall offset is ticket offset * user limit + user offset
	offset := ticketOffset*userLimit + userOffset
	match, err := c.service.database.GetMatch(ctx, model.GetMatchParams{
		Match: model.MatchParams{
			MatchmakingTicket: model.MatchmakingTicketParams{
				MatchmakingUser: model.GetMatchmakingUserParams{
					ID:           conversion.Uint64ToSqlNullInt64(c.In.Match.MatchmakingTicket.Id),
					ClientUserID: conversion.Uint64ToSqlNullInt64(c.In.Match.MatchmakingTicket.MatchmakingUser.ClientUserId),
				},
				ID: conversion.Uint64ToSqlNullInt64(c.In.Match.MatchmakingTicket.Id),
			},
			ID: conversion.Uint64ToSqlNullInt64(c.In.Match.Id),
		},
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			c.Out = &api.GetMatchResponse{
				Success: false,
				Error:   api.GetMatchResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	apiMatch, err := unmarshalMatch(match)
	if err != nil {
		return err
	}
	c.Out = &api.GetMatchResponse{
		Success: true,
		Match:   apiMatch,
		Error:   api.GetMatchResponse_NONE,
	}
	return nil
}
