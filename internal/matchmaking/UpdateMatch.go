package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type UpdateMatchCommand struct {
	service *Service
	In      *api.UpdateMatchRequest
	Out     *api.UpdateMatchResponse
}

func NewUpdateMatchCommand(service *Service, in *api.UpdateMatchRequest) *UpdateMatchCommand {
	return &UpdateMatchCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateMatchCommand) Execute(ctx context.Context) error {
	mmErr := c.service.checkForMatchRequestError(c.In.Match)
	// Check if error is found
	if mmErr != nil {
		c.Out = &api.UpdateMatchResponse{
			Success: false,
			Error:   conversion.Enum(*mmErr, api.UpdateMatchResponse_Error_value, api.UpdateMatchResponse_NOT_FOUND),
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
	// Check if data is given
	if c.In.Data == nil {
		c.Out = &api.UpdateMatchResponse{
			Success: false,
			Error:   api.UpdateMatchResponse_DATA_REQUIRED,
		}
		return nil
	}
	// Prepare data
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	result, err := c.service.database.UpdateMatch(ctx, model.UpdateMatchParams{
		Match: model.MatchParams{
			MatchmakingTicket: model.MatchmakingTicketParams{
				MatchmakingUser: model.MatchmakingUserParams{
					ID:           conversion.Uint64ToSqlNullInt64(c.In.Match.MatchmakingTicket.Id),
					ClientUserID: conversion.Uint64ToSqlNullInt64(c.In.Match.MatchmakingTicket.MatchmakingUser.ClientUserId),
				},
				ID:       conversion.Uint64ToSqlNullInt64(c.In.Match.MatchmakingTicket.Id),
				Statuses: []string{"PENDING", "MATCHED"},
			},
			ID: conversion.Uint64ToSqlNullInt64(c.In.Match.Id),
		},
		Data: data,
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
		match, err := c.service.database.GetMatch(ctx, model.GetMatchParams{
			Match: model.MatchParams{
				MatchmakingTicket: model.MatchmakingTicketParams{
					MatchmakingUser: model.MatchmakingUserParams{
						ID:           conversion.Uint64ToSqlNullInt64(c.In.Match.MatchmakingTicket.Id),
						ClientUserID: conversion.Uint64ToSqlNullInt64(c.In.Match.MatchmakingTicket.MatchmakingUser.ClientUserId),
					},
					ID:       conversion.Uint64ToSqlNullInt64(c.In.Match.MatchmakingTicket.Id),
					Statuses: []string{"PENDING", "MATCHED"},
				},
				ID: conversion.Uint64ToSqlNullInt64(c.In.Match.Id),
			},
			TicketLimit: 1,
			UserLimit:   1,
			ArenaLimit:  1,
		})
		if err != nil {
			return err
		}
		if len(match) == 0 {
			c.Out = &api.UpdateMatchResponse{
				Success: false,
				Error:   api.UpdateMatchResponse_NOT_FOUND,
			}
			return nil
		}
	}
	c.Out = &api.UpdateMatchResponse{
		Success: true,
		Error:   api.UpdateMatchResponse_NONE,
	}
	return nil
}
