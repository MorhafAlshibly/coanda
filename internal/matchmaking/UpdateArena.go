package matchmaking

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"open-match.dev/open-match/pkg/pb"
)

type UpdateArenaCommand struct {
	service *Service
	In      *api.UpdateArenaRequest
	Out     *api.UpdateArenaResponse
}

func NewUpdateArenaCommand(service *Service, in *api.UpdateArenaRequest) *UpdateArenaCommand {
	return &UpdateArenaCommand{
		service: service,
		In:      in,
	}
}

func (c *UpdateArenaCommand) Execute(ctx context.Context) error {
	aErr := c.service.checkForArenaRequestError(c.In.Arena)
	// Check if error is found
	if aErr != nil {
		c.Out = &api.UpdateArenaResponse{
			Success: false,
			Error:   conversion.Enum(*aErr, api.UpdateArenaResponse_Error_value, api.UpdateArenaResponse_NOT_FOUND),
		}
		return nil
	}
	// Check if no update is specified
	if c.In.MinPlayers == nil && c.In.MaxPlayers == nil && c.In.Data == nil {
		c.Out = &api.UpdateArenaResponse{
			Success: false,
			Error:   api.UpdateArenaResponse_NO_UPDATE_SPECIFIED,
		}
		return nil
	}
	// Prepare data
	data := json.RawMessage(nil)
	if c.In.Data != nil {
		var err error
		data, err = conversion.ProtobufStructToRawJson(c.In.Data)
		if err != nil {
			return err
		}
	}
	if c.In.MinPlayers != nil || c.In.MaxPlayersPerTicket != nil || c.In.MaxPlayers != nil {
		if c.In.MinPlayers == nil || c.In.MaxPlayers == nil || c.In.MaxPlayersPerTicket == nil {
			c.Out = &api.UpdateArenaResponse{
				Success: false,
				Error:   api.UpdateArenaResponse_IF_CAPACITY_CHANGED_MUST_CHANGE_ALL_PLAYERS,
			}
			return nil
		}
		if *c.In.MinPlayers > *c.In.MaxPlayers {
			c.Out = &api.UpdateArenaResponse{
				Success: false,
				Error:   api.UpdateArenaResponse_MIN_PLAYERS_CANNOT_BE_GREATER_THAN_MAX_PLAYERS,
			}
			return nil
		}
		if *c.In.MaxPlayersPerTicket < *c.In.MinPlayers {
			c.Out = &api.UpdateArenaResponse{
				Success: false,
				Error:   api.UpdateArenaResponse_MAX_PLAYERS_PER_TICKET_CANNOT_BE_LESS_THAN_MIN_PLAYERS,
			}
			return nil
		}
		if *c.In.MaxPlayersPerTicket > *c.In.MaxPlayers {
			c.Out = &api.UpdateArenaResponse{
				Success: false,
				Error:   api.UpdateArenaResponse_MAX_PLAYERS_PER_TICKET_CANNOT_BE_GREATER_THAN_MAX_PLAYERS,
			}
			return nil
		}
	}
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	arena, err := qtx.GetArena(ctx, model.ArenaParams{
		ID:   conversion.Uint64ToSqlNullInt64(c.In.Arena.Id),
		Name: conversion.StringToSqlNullString(c.In.Arena.Name),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			c.Out = &api.UpdateArenaResponse{
				Success: false,
				Error:   api.UpdateArenaResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	// Get any tickets that are currently queuing the arena
	ticketClient, err := c.service.queryServiceClient.QueryTicketIds(ctx, &pb.QueryTicketIdsRequest{
		Pool: &pb.Pool{
			Name: "default",
			TagPresentFilters: []*pb.TagPresentFilter{
				{
					Tag: fmt.Sprintf("Arena_%d", arena.ID),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = ticketClient.Recv()
	if err != io.EOF {
		c.Out = &api.UpdateArenaResponse{
			Success: false,
			Error:   api.UpdateArenaResponse_ARENA_CURRENTLY_IN_USE,
		}
		return nil
	}
	_, err = qtx.UpdateArena(ctx, model.UpdateArenaParams{
		Arena: model.ArenaParams{
			ID:   conversion.Uint64ToSqlNullInt64(c.In.Arena.Id),
			Name: conversion.StringToSqlNullString(c.In.Arena.Name),
		},
		Data:                data,
		MinPlayers:          conversion.Uint32ToSqlNullInt32(c.In.MinPlayers),
		MaxPlayersPerTicket: conversion.Uint32ToSqlNullInt32(c.In.MaxPlayersPerTicket),
		MaxPlayers:          conversion.Uint32ToSqlNullInt32(c.In.MaxPlayers),
	})
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.UpdateArenaResponse{
		Success: true,
		Error:   api.UpdateArenaResponse_NONE,
	}
	return nil
}
