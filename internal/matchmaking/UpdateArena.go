package matchmaking

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/goquOptions"
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
			Error:   conversion.Enum(*aErr, api.UpdateArenaResponse_Error_value, api.UpdateArenaResponse_ARENA_ID_OR_NAME_REQUIRED),
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
	// Check if the arena exists and lock it for update
	_, err = qtx.GetArena(ctx, arenaRequestToArenaParams(c.In.Arena), &goquOptions.SelectDataset{Locked: true})
	if err != nil {
		if err == sql.ErrNoRows {
			// If we didn't find a row
			c.Out = &api.UpdateArenaResponse{
				Success: false,
				Error:   api.UpdateArenaResponse_NOT_FOUND,
			}
			return nil
		}
		return err
	}
	// Get any tickets that are currently queuing the arena
	tickets, err := qtx.GetMatchmakingTickets(ctx, model.GetMatchmakingTicketsParams{
		Arena:    arenaRequestToArenaParams(c.In.Arena),
		Statuses: []string{"PENDING", "MATCHED"},
	})
	if err != nil {
		return err
	}
	if len(tickets) > 0 {
		c.Out = &api.UpdateArenaResponse{
			Success: false,
			Error:   api.UpdateArenaResponse_ARENA_CURRENTLY_IN_USE,
		}
		return nil
	}
	result, err := qtx.UpdateArena(ctx, model.UpdateArenaParams{
		Arena:               arenaRequestToArenaParams(c.In.Arena),
		Data:                data,
		MinPlayers:          conversion.Uint32ToSqlNullInt32(c.In.MinPlayers),
		MaxPlayersPerTicket: conversion.Uint32ToSqlNullInt32(c.In.MaxPlayersPerTicket),
		MaxPlayers:          conversion.Uint32ToSqlNullInt32(c.In.MaxPlayers),
	})
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
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
