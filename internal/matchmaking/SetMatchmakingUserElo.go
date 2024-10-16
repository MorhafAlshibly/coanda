package matchmaking

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/errorcode"
	"github.com/go-sql-driver/mysql"
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
			Error:   conversion.Enum(*muErr, api.SetMatchmakingUserEloResponse_Error_value, api.SetMatchmakingUserEloResponse_MATCHMAKING_USER_ID_OR_CLIENT_USER_ID_REQUIRED),
		}
		return nil
	}
	maErr := c.service.checkForArenaRequestError(c.In.Arena)
	// Check if error is found
	if maErr != nil {
		c.Out = &api.SetMatchmakingUserEloResponse{
			Success: false,
			Error:   conversion.Enum(*maErr, api.SetMatchmakingUserEloResponse_Error_value, api.SetMatchmakingUserEloResponse_ID_OR_NAME_REQUIRED),
		}
		return nil
	}
	// Start transaction
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	user, err := qtx.GetMatchmakingUser(ctx, model.GetMatchmakingUserParams{
		ID:           conversion.Uint64ToSqlNullInt64(c.In.MatchmakingUser.Id),
		ClientUserID: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingUser.ClientUserId),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			c.Out = &api.SetMatchmakingUserEloResponse{
				Success: false,
				Error:   api.SetMatchmakingUserEloResponse_USER_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	arena, err := qtx.GetArena(ctx, model.GetArenaParams{
		ID:   conversion.Uint64ToSqlNullInt64(c.In.Arena.Id),
		Name: conversion.StringToSqlNullString(c.In.Arena.Name),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			c.Out = &api.SetMatchmakingUserEloResponse{
				Success: false,
				Error:   api.SetMatchmakingUserEloResponse_ARENA_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	// Create user elo
	result, err := qtx.CreateMatchmakingUserElo(ctx, model.CreateMatchmakingUserEloParams{
		MatchmakingUserID:  user.ID,
		MatchmakingArenaID: arena.ID,
		Elo:                int32(c.In.Elo),
	})
	if err != nil {
		// If the elo already exists, update it
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if errorcode.IsDuplicateEntry(mysqlErr, "matchmaking_user_elo", "matchmaking_user_elo_matchmaking_user_id_matchmaking_arena_id") {
				result, err = qtx.UpdateMatchmakingUserElo(ctx, model.UpdateMatchmakingUserEloParams{
					MatchmakingUserElo: model.GetMatchmakingUserEloParams{
						MatchmakingUser: model.GetMatchmakingUserParams{
							ID: conversion.Uint64ToSqlNullInt64(&user.ID),
						},
						Arena: model.GetArenaParams{
							ID: conversion.Uint64ToSqlNullInt64(&arena.ID),
						},
					},
					Elo:          int32(c.In.Elo),
					IncrementElo: c.In.IncrementElo,
				})
				if err != nil {
					return err
				}
				rowsAffected, err := result.RowsAffected()
				if err != nil {
					return err
				}
				if rowsAffected == 0 {
					// Unexpected error
					return errors.New("unable to update elo")
				}
			} else {
				return err
			}
		} else {
			return err
		}
	}
	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.SetMatchmakingUserEloResponse{
		Success: true,
		Error:   api.SetMatchmakingUserEloResponse_NONE,
	}
	return nil
}
