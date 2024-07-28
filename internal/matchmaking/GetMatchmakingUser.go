package matchmaking

import (
	"context"
	"database/sql"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetMatchmakingUserCommand struct {
	service *Service
	In      *api.MatchmakingUserRequest
	Out     *api.GetMatchmakingUserResponse
}

func NewGetMatchmakingUserCommand(service *Service, in *api.MatchmakingUserRequest) *GetMatchmakingUserCommand {
	return &GetMatchmakingUserCommand{
		service: service,
		In:      in,
	}
}

func (c *GetMatchmakingUserCommand) Execute(ctx context.Context) error {
	muErr := c.service.checkForMatchmakingUserRequestError(c.In)
	if muErr != nil {
		c.Out = &api.GetMatchmakingUserResponse{
			Success: false,
			Error:   conversion.Enum(*muErr, api.GetMatchmakingUserResponse_Error_value, api.GetMatchmakingUserResponse_MATCHMAKING_USER_ID_OR_CLIENT_USER_ID_REQUIRED),
		}
		return nil
	}
	matchmakingUser, err := c.service.database.GetMatchmakingUser(ctx, model.GetMatchmakingUserParams{
		ID:           conversion.Uint64ToSqlNullInt64(c.In.Id),
		ClientUserID: conversion.Uint64ToSqlNullInt64(c.In.ClientUserId),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			c.Out = &api.GetMatchmakingUserResponse{
				Success: false,
				Error:   api.GetMatchmakingUserResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	apiMatchmakingUser, err := unmarshalMatchmakingUser(matchmakingUser)
	if err != nil {
		return err
	}
	c.Out = &api.GetMatchmakingUserResponse{
		Success:         true,
		MatchmakingUser: apiMatchmakingUser,
	}
	return nil
}
