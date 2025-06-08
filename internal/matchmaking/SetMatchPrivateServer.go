package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type SetMatchPrivateServerCommand struct {
	service *Service
	In      *api.SetMatchPrivateServerRequest
	Out     *api.SetMatchPrivateServerResponse
}

func NewSetMatchPrivateServerCommand(service *Service, in *api.SetMatchPrivateServerRequest) *SetMatchPrivateServerCommand {
	return &SetMatchPrivateServerCommand{
		service: service,
		In:      in,
	}
}

func (c *SetMatchPrivateServerCommand) Execute(ctx context.Context) error {
	mmErr := c.service.checkForMatchRequestError(c.In.Match)
	// Check if error is found
	if mmErr != nil {
		c.Out = &api.SetMatchPrivateServerResponse{
			Success: false,
			Error:   conversion.Enum(*mmErr, api.SetMatchPrivateServerResponse_Error_value, api.SetMatchPrivateServerResponse_MATCH_ID_OR_MATCHMAKING_TICKET_REQUIRED),
		}
		return nil
	}
	// Check if private server id is given
	if c.In.PrivateServerId == "" {
		c.Out = &api.SetMatchPrivateServerResponse{
			Success: false,
			Error:   api.SetMatchPrivateServerResponse_PRIVATE_SERVER_ID_REQUIRED,
		}
		return nil
	}
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := model.New(tx)
	params := matchRequestToMatchParams(c.In.Match)
	result, err := qtx.SetMatchPrivateServer(ctx, model.SetMatchPrivateServerParams{
		Match:           params,
		PrivateServerID: c.In.PrivateServerId,
	})
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 0 {
		err = tx.Commit()
		if err != nil {
			return err
		}
		c.Out = &api.SetMatchPrivateServerResponse{
			Success:         true,
			PrivateServerId: &c.In.PrivateServerId,
			Error:           api.SetMatchPrivateServerResponse_NONE,
		}
		return nil
	}
	match, err := qtx.GetMatch(ctx, model.GetMatchParams{
		Match:       params,
		TicketLimit: 1,
		UserLimit:   1,
		ArenaLimit:  1,
	})
	if err != nil {
		return err
	}
	if len(match) == 0 {
		c.Out = &api.SetMatchPrivateServerResponse{
			Success: false,
			Error:   api.SetMatchPrivateServerResponse_NOT_FOUND,
		}
		return nil
	}
	c.Out = &api.SetMatchPrivateServerResponse{
		Success:         false,
		PrivateServerId: conversion.SqlNullStringToString(match[0].PrivateServerID),
		Error:           api.SetMatchPrivateServerResponse_PRIVATE_SERVER_ALREADY_SET,
	}
	return nil
}
