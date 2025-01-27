package matchmaking

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type ExpireMatchmakingTicketCommand struct {
	service *Service
	In      *api.MatchmakingTicketRequest
	Out     *api.ExpireMatchmakingTicketResponse
}

func NewExpireMatchmakingTicketCommand(service *Service, in *api.MatchmakingTicketRequest) *ExpireMatchmakingTicketCommand {
	return &ExpireMatchmakingTicketCommand{
		service: service,
		In:      in,
	}
}

func (c *ExpireMatchmakingTicketCommand) Execute(ctx context.Context) error {
	mtErr := c.service.checkForMatchmakingTicketRequestError(c.In)
	if mtErr != nil {
		c.Out = &api.ExpireMatchmakingTicketResponse{
			Success: false,
			Error:   conversion.Enum(*mtErr, api.ExpireMatchmakingTicketResponse_Error_value, api.ExpireMatchmakingTicketResponse_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED),
		}
		return nil
	}
	// Make sure matchmaking user isnt nil
	if c.In.MatchmakingUser == nil {
		c.In.MatchmakingUser = &api.MatchmakingUserRequest{}
	}
	result, err := c.service.database.ExpireMatchmakingTicket(ctx, model.MatchmakingTicketParams{
		MatchmakingUser: model.MatchmakingUserParams{
			ID:           conversion.Uint64ToSqlNullInt64(c.In.Id),
			ClientUserID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingUser.ClientUserId),
		},
		ID:       conversion.Uint64ToSqlNullInt64(c.In.Id),
		Statuses: []string{"PENDING"},
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
		ticket, err := c.service.database.GetMatchmakingTicket(ctx, model.GetMatchmakingTicketParams{
			MatchmakingTicket: model.MatchmakingTicketParams{
				MatchmakingUser: model.MatchmakingUserParams{
					ID:           conversion.Uint64ToSqlNullInt64(c.In.Id),
					ClientUserID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingUser.ClientUserId),
				},
				ID:                        conversion.Uint64ToSqlNullInt64(c.In.Id),
				Statuses:                  []string{"PENDING", "MATCHED"},
				GetByIDRegardlessOfStatus: true,
			},
			UserLimit:  1,
			ArenaLimit: 1,
		})
		// Check if ticket is found
		if err != nil {
			return err
		}
		if len(ticket) == 0 {
			c.Out = &api.ExpireMatchmakingTicketResponse{
				Success: false,
				Error:   api.ExpireMatchmakingTicketResponse_NOT_FOUND,
			}
			return nil
		}
		if ticket[0].Status == "EXPIRED" {
			c.Out = &api.ExpireMatchmakingTicketResponse{
				Success: false,
				Error:   api.ExpireMatchmakingTicketResponse_ALREADY_EXPIRED,
			}
			return nil
		}
		// If ticket is matched we can't expire it
		if ticket[0].Status == "MATCHED" {
			c.Out = &api.ExpireMatchmakingTicketResponse{
				Success: false,
				Error:   api.ExpireMatchmakingTicketResponse_ALREADY_MATCHED,
			}
			return nil
		}
		if ticket[0].Status == "ENDED" {
			c.Out = &api.ExpireMatchmakingTicketResponse{
				Success: false,
				Error:   api.ExpireMatchmakingTicketResponse_ALREADY_ENDED,
			}
			return nil
		}
		if ticket[0].Status == "PENDING" {
			// Unexpected error
			return errors.New("could not expire ticket")
		}
	}
	c.Out = &api.ExpireMatchmakingTicketResponse{
		Success: true,
		Error:   api.ExpireMatchmakingTicketResponse_NONE,
	}
	return nil
}
