package matchmaking

import (
	"context"
	"database/sql"
	"time"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type CreateMatchmakingTicketCommand struct {
	service *Service
	In      *api.CreateMatchmakingTicketRequest
	Out     *api.CreateMatchmakingTicketResponse
}

func NewCreateMatchmakingTicketCommand(service *Service, in *api.CreateMatchmakingTicketRequest) *CreateMatchmakingTicketCommand {
	return &CreateMatchmakingTicketCommand{
		service: service,
		In:      in,
	}
}

func (c *CreateMatchmakingTicketCommand) Execute(ctx context.Context) error {
	if len(c.In.MatchmakingUsers) == 0 {
		c.Out = &api.CreateMatchmakingTicketResponse{
			Success: false,
			Error:   api.CreateMatchmakingTicketResponse_MATCHMAKING_USERS_REQUIRED,
		}
		return nil
	}
	if len(c.In.Arenas) == 0 {
		c.Out = &api.CreateMatchmakingTicketResponse{
			Success: false,
			Error:   api.CreateMatchmakingTicketResponse_ARENAS_REQUIRED,
		}
		return nil
	}
	if c.In.Data == nil {
		c.Out = &api.CreateMatchmakingTicketResponse{
			Success: false,
			Error:   api.CreateMatchmakingTicketResponse_DATA_REQUIRED,
		}
		return nil
	}
	data, err := conversion.ProtobufStructToRawJson(c.In.Data)
	if err != nil {
		return err
	}
	tx, err := c.service.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := c.service.database.WithTx(tx)
	// Get all arena ids
	numberOfUsers := uint32(len(c.In.MatchmakingUsers))
	arenaIds := make([]uint64, 0, len(c.In.Arenas))
	for _, arena := range c.In.Arenas {
		arena, err := qtx.GetArena(ctx, model.GetArenaParams{
			ID:   conversion.Uint64ToSqlNullInt64(arena.Id),
			Name: conversion.StringToSqlNullString(arena.Name),
		})
		if err != nil {
			if err == sql.ErrNoRows {
				c.Out = &api.CreateMatchmakingTicketResponse{
					Success: false,
					Error:   api.CreateMatchmakingTicketResponse_ARENA_NOT_FOUND,
				}
				return nil
			}
			return err
		}
		// Check if too many players
		if numberOfUsers > uint32(arena.MaxPlayersPerTicket) {
			c.Out = &api.CreateMatchmakingTicketResponse{
				Success: false,
				Error:   api.CreateMatchmakingTicketResponse_TOO_MANY_PLAYERS,
			}
			return nil
		}
		arenaIds = append(arenaIds, arena.ID)
	}
	// Get all user ids
	userIds := make([]uint64, 0, len(c.In.MatchmakingUsers))
	for _, user := range c.In.MatchmakingUsers {
		user, err := qtx.GetMatchmakingUser(ctx, model.MatchmakingUserParams{
			ID:           conversion.Uint64ToSqlNullInt64(user.Id),
			ClientUserID: conversion.Uint64ToSqlNullInt64(user.ClientUserId),
		})
		if err != nil {
			if err == sql.ErrNoRows {
				c.Out = &api.CreateMatchmakingTicketResponse{
					Success: false,
					Error:   api.CreateMatchmakingTicketResponse_USER_NOT_FOUND,
				}
				return nil
			}
			return err
		}
		// Check if user has an active ticket
		ticket, err := qtx.GetMatchmakingTicket(ctx, model.GetMatchmakingTicketParams{
			MatchmakingTicket: model.MatchmakingTicketParams{
				MatchmakingUser: model.MatchmakingUserParams{
					ID: conversion.Uint64ToSqlNullInt64(&user.ID),
				},
				Statuses: []string{"PENDING", "MATCHED"},
			},
			UserLimit:  1,
			ArenaLimit: 1,
		})
		if err != nil {
			return err
		}
		if len(ticket) > 0 {
			c.Out = &api.CreateMatchmakingTicketResponse{
				Success: false,
				Error:   api.CreateMatchmakingTicketResponse_USER_ALREADY_HAS_ACTIVE_TICKET,
			}
			return nil
		}
		userIds = append(userIds, user.ID)
	}
	// Create the ticket
	result, err := qtx.CreateMatchmakingTicket(ctx, model.CreateMatchmakingTicketParams{
		Data:      data,
		EloWindow: 0,
		ExpiresAt: time.Now().Add(c.service.expiryTimeWindow),
	})
	if err != nil {
		return err
	}
	ticketId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	// Add the users to the ticket
	for _, userId := range userIds {
		_, err := qtx.CreateMatchmakingTicketUser(ctx, model.CreateMatchmakingTicketUserParams{
			MatchmakingTicketID: uint64(ticketId),
			MatchmakingUserID:   userId,
		})
		if err != nil {
			return err
		}
	}
	// Add the arenas to the ticket
	for _, arenaId := range arenaIds {
		_, err := qtx.CreateMatchmakingTicketArena(ctx, model.CreateMatchmakingTicketArenaParams{
			MatchmakingTicketID: uint64(ticketId),
			MatchmakingArenaID:  arenaId,
		})
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	c.Out = &api.CreateMatchmakingTicketResponse{
		Success: true,
		Id:      conversion.ValueToPointer(uint64(ticketId)),
		Error:   api.CreateMatchmakingTicketResponse_NONE,
	}
	return nil
}
