package matchmaking

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type UpdateMatchmakingTicketCommand struct {
	service *Service
	In      *api.UpdateMatchmakingTicketRequest
	Out     *api.UpdateMatchmakingTicketResponse
}

func NewUpdateMatchmakingTicketCommand(service *Service, in *api.UpdateMatchmakingTicketRequest) *UpdateMatchmakingTicketCommand {
	return &UpdateMatchmakingTicketCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateMatchmakingTicketCommand) Execute(ctx context.Context) error {
	mtErr := c.service.checkForMatchmakingTicketRequestError(c.In.MatchmakingTicket)
	// Check if error is found
	if mtErr != nil {
		c.Out = &api.UpdateMatchmakingTicketResponse{
			Success: false,
			Error:   conversion.Enum(*mtErr, api.UpdateMatchmakingTicketResponse_Error_value, api.UpdateMatchmakingTicketResponse_NOT_FOUND),
		}
		return nil
	}
	// Make sure matchmaking user isnt nil
	if c.In.MatchmakingTicket.MatchmakingUser == nil {
		c.In.MatchmakingTicket.MatchmakingUser = &api.MatchmakingUserRequest{}
	}
	// Check if data is given
	if c.In.Data == nil {
		c.Out = &api.UpdateMatchmakingTicketResponse{
			Success: false,
			Error:   api.UpdateMatchmakingTicketResponse_DATA_REQUIRED,
		}
		return nil
	}
	// Prepare data
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	result, err := c.service.database.UpdateMatchmakingTicket(ctx, model.UpdateMatchmakingTicketParams{
		MatchmakingTicket: model.MatchmakingTicketParams{
			ID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.Id),
			MatchmakingUser: model.MatchmakingUserParams{
				ID:           conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.MatchmakingUser.Id),
				ClientUserID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.MatchmakingUser.ClientUserId),
			},
			Statuses: []string{"PENDING", "MATCHED"},
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
		_, err = c.service.database.GetMatchmakingTicket(ctx, model.GetMatchmakingTicketParams{
			MatchmakingTicket: model.MatchmakingTicketParams{
				ID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.Id),
				MatchmakingUser: model.MatchmakingUserParams{
					ID:           conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.MatchmakingUser.Id),
					ClientUserID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.MatchmakingUser.ClientUserId),
				},
				Statuses: []string{"PENDING", "MATCHED"},
			},
			UserLimit:  1,
			ArenaLimit: 1,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				// If we didn't find a row
				c.Out = &api.UpdateMatchmakingTicketResponse{
					Success: false,
					Error:   api.UpdateMatchmakingTicketResponse_NOT_FOUND,
				}
				return nil
			}
			return err
		}
	}
	c.Out = &api.UpdateMatchmakingTicketResponse{
		Success: true,
		Error:   api.UpdateMatchmakingTicketResponse_NONE,
	}
	return nil
}
