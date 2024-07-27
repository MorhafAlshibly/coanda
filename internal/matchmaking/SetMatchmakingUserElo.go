package matchmaking

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type SetMatchmakingUserEloCommand struct {
	service *Service
	In      *api.SetMatchmakingUserEloRequest
	Out     *api.SetMatchmakingUserEloResponse
}

func NewSetMatchmakingUserEloCommand(service *Service, in *api.SetMatchmakingUserEloRequest) *SetMatchmakingUserEloCommand {
	return &SetMatchmakingUserEloCommand{
		service: service,
		In:      in,
	}
}

func (c *SetMatchmakingUserEloCommand) Execute(ctx context.Context) error {
	muErr := c.service.checkForMatchmakingUserRequestError(c.In.MatchmakingUser)
	// Check if error is found
	if muErr != nil {
		c.Out = &api.SetMatchmakingUserEloResponse{
			Success: false,
			Error:   conversion.Enum(*muErr, api.SetMatchmakingUserEloResponse_Error_value, api.SetMatchmakingUserEloResponse_NOT_FOUND),
		}
		return nil
	}
	// Check if elo is given
	if c.In.Elo == nil {
		c.Out = &api.SetMatchmakingUserEloResponse{
			Success: false,
			Error:   api.SetMatchmakingUserEloResponse_ELO_REQUIRED,
		}
		return nil
	}
	result, err := c.service.database.SetMatchmakingUserElo(ctx, model.SetMatchmakingUserEloParams{
		MatchmakingUser: model.GetMatchmakingUserParams{
			ID:     conversion.Uint64ToSqlNullInt64(c.In.MatchmakingUser.Id),
			UserID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingUser.UserId),
		},
		Elo: conversion.Int64ToSqlNullInt64(c.In.Elo),
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
		_, err = c.service.database.GetMatchmakingUser(ctx, model.GetMatchmakingUserParams{
			ID:     conversion.Uint64ToSqlNullInt64(c.In.MatchmakingUser.Id),
			UserID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingUser.UserId),
		})
		if err != nil {
			if err == sql.ErrNoRows {
				// If we didn't find a row
				c.Out = &api.SetMatchmakingUserEloResponse{
					Success: false,
					Error:   api.SetMatchmakingUserEloResponse_NOT_FOUND,
				}
				return nil
			}
			return err
		}
	}
	c.Out = &api.SetMatchmakingUserEloResponse{
		Success: true,
		Error:   api.SetMatchmakingUserEloResponse_NONE,
	}
	return nil
}
