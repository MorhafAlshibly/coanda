package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type DeleteMatchCommand struct {
	service *Service
	In      *api.MatchRequest
	Out     *api.DeleteMatchResponse
}

func NewDeleteMatchCommand(service *Service, in *api.MatchRequest) *DeleteMatchCommand {
	return &DeleteMatchCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteMatchCommand) Execute(ctx context.Context) error {
	mmErr := c.service.checkForMatchRequestError(c.In)
	// Check if error is found
	if mmErr != nil {
		c.Out = &api.DeleteMatchResponse{
			Success: false,
			Error:   conversion.Enum(*mmErr, api.DeleteMatchResponse_Error_value, api.DeleteMatchResponse_NOT_FOUND),
		}
		return nil
	}
	// Make sure matchmaking ticket isnt nil
	if c.In.MatchmakingTicket == nil {
		c.In.MatchmakingTicket = &api.MatchmakingTicketRequest{
			MatchmakingUser: &api.MatchmakingUserRequest{},
		}
	}
	// Make sure matchmaking user isnt nil
	if c.In.MatchmakingTicket.MatchmakingUser == nil {
		c.In.MatchmakingTicket.MatchmakingUser = &api.MatchmakingUserRequest{}
	}
	params := model.MatchParams{
		MatchmakingTicket: model.MatchmakingTicketParams{
			MatchmakingUser: model.MatchmakingUserParams{
				ID:           conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.Id),
				ClientUserID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.MatchmakingUser.ClientUserId),
			},
			ID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.Id),
		},
		ID: conversion.Uint64ToSqlNullInt64(c.In.Id),
	}
	result, err := c.service.database.DeleteMatch(ctx, params)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// Check if we didn't find a row
		match, err := c.service.database.GetMatch(ctx, model.GetMatchParams{
			Match:       params,
			TicketLimit: 1,
			UserLimit:   1,
			ArenaLimit:  1,
		})
		if err != nil {
			return err
		}
		if len(match) == 0 {
			c.Out = &api.DeleteMatchResponse{
				Success: false,
				Error:   api.DeleteMatchResponse_NOT_FOUND,
			}
			return nil
		}
	}
	c.Out = &api.DeleteMatchResponse{
		Success: true,
		Error:   api.DeleteMatchResponse_NONE,
	}
	return nil
}
