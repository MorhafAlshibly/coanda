package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type PollMatchmakingTicketCommand struct {
	service *Service
	In      *api.MatchmakingTicketRequest
	Out     *api.MatchmakingTicketResponse
}

func NewPollMatchmakingTicketCommand(service *Service, in *api.MatchmakingTicketRequest) *PollMatchmakingTicketCommand {
	return &PollMatchmakingTicketCommand{
		service: service,
		In:      in,
	}
}

func (c *PollMatchmakingTicketCommand) Execute(ctx context.Context) error {
	mtErr := c.service.checkForMatchmakingTicketRequestError(c.In)
	if mtErr != nil {
		c.Out = &api.MatchmakingTicketResponse{
			Success: false,
			Error:   conversion.Enum(*mtErr, api.MatchmakingTicketResponse_Error_value, api.MatchmakingTicketResponse_ID_OR_MATCHMAKING_USER_REQUIRED),
		}
		return nil
	}
	// Make sure matchmaking user isnt nil
	if c.In.MatchmakingUser == nil {
		c.In.MatchmakingUser = &api.MatchmakingUserRequest{}
	}
	result, err := c.service.database.PollMatchmakingTicket(ctx, model.PollMatchmakingTicketParams{
		MatchmakingUser: model.GetMatchmakingUserParams{
			ID:           conversion.Uint64ToSqlNullInt64(c.In.Id),
			ClientUserID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingUser.ClientUserId),
		},
		ID: conversion.Uint64ToSqlNullInt64(c.In.Id),
	})
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// Check if we didn't find a row
		c.Out = &api.MatchmakingTicketResponse{
			Success: false,
			Error:   api.MatchmakingTicketResponse_NOT_FOUND,
		}
		return nil
	}
	c.Out = &api.MatchmakingTicketResponse{
		Success: true,
		Error:   api.MatchmakingTicketResponse_NONE,
	}
	return nil
}
