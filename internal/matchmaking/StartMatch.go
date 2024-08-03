package matchmaking

import (
	"context"
	"database/sql"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type StartMatchCommand struct {
	service *Service
	In      *api.StartMatchRequest
	Out     *api.StartMatchResponse
}

func NewStartMatchCommand(service *Service, in *api.StartMatchRequest) *StartMatchCommand {
	return &StartMatchCommand{
		service: service,
		In:      in,
	}
}

func (c *StartMatchCommand) Execute(ctx context.Context) error {
	mmErr := c.service.checkForMatchRequestError(c.In.Match)
	// Check if error is found
	if mmErr != nil {
		c.Out = &api.StartMatchResponse{
			Success: false,
			Error:   conversion.Enum(*mmErr, api.StartMatchResponse_Error_value, api.StartMatchResponse_NOT_FOUND),
		}
		return nil
	}
	// Check if start time is nil
	if c.In.StartTime == nil {
		c.Out = &api.StartMatchResponse{
			Success: false,
			Error:   api.StartMatchResponse_START_TIME_REQUIRED,
		}
		return nil
	}
	if c.In.StartTime.AsTime().Before(time.Now()) {
		c.Out = &api.StartMatchResponse{
			Success: false,
			Error:   api.StartMatchResponse_INVALID_START_TIME,
		}
		return nil
	}
	if c.In.StartTime.AsTime().Before(time.Now().Add(c.service.startTimeBuffer)) {
		c.Out = &api.StartMatchResponse{
			Success: false,
			Error:   api.StartMatchResponse_START_TIME_TOO_SOON,
		}
		return nil
	}
	lockTime := c.In.StartTime.AsTime().Add(-c.service.lockedAtBuffer)
	// Check if lock time is before now
	if lockTime.Before(time.Now()) {
		c.Out = &api.StartMatchResponse{
			Success: false,
			Error:   api.StartMatchResponse_START_TIME_TOO_SOON,
		}
		return nil
	}
	// Start match
	result, err := c.service.database.StartMatch(ctx, model.StartMatchParams{
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
		StartTime: c.In.StartTime.AsTime(),
		LockTime:  lockTime,
	})
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// Either the match wasnt found or it has already started
		_, err := c.service.database.GetMatch(ctx, model.GetMatchParams{
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
			Limit:  1,
			Offset: 0,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				// If we didn't find a row
				c.Out = &api.StartMatchResponse{
					Success: false,
					Error:   api.StartMatchResponse_NOT_FOUND,
				}
				return nil
			}
			return err
		}
		c.Out = &api.StartMatchResponse{
			Success: false,
			Error:   api.StartMatchResponse_ALREADY_STARTED,
		}
		return nil
	}
	c.Out = &api.StartMatchResponse{
		Success: true,
		Error:   api.StartMatchResponse_NONE,
	}
	return nil
}
