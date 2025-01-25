package matchmaking

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"strconv"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"google.golang.org/protobuf/types/known/anypb"
	"open-match.dev/open-match/pkg/pb"
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
	data, err := anypb.New(c.In.Data)
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
	arenaTags := make([]string, 0, len(c.In.Arenas))
	floatTicketMap := map[string]float64{}
	for _, arenaRequest := range c.In.Arenas {
		arena, err := qtx.GetArena(ctx, model.ArenaParams{
			ID:   conversion.Uint64ToSqlNullInt64(arenaRequest.Id),
			Name: conversion.StringToSqlNullString(arenaRequest.Name),
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
		arenaTag := fmt.Sprintf("Arena_%d", arena.ID)
		floatTicketMap[arenaTag] = float64(arena.MaxPlayers)
		arenaTags = append(arenaTags, arenaTag)
	}
	// Get all user ids
	userTags := make([]string, 0, len(c.In.MatchmakingUsers))
	for _, userRequest := range c.In.MatchmakingUsers {
		user, err := qtx.GetMatchmakingUser(ctx, model.MatchmakingUserParams{
			ID:           conversion.Uint64ToSqlNullInt64(userRequest.Id),
			ClientUserID: conversion.Uint64ToSqlNullInt64(userRequest.ClientUserId),
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
		userTag := fmt.Sprintf("User_%d", user.ClientUserID)
		// Check if user has an active ticket
		ticketClient, err := c.service.queryServiceClient.QueryTicketIds(ctx, &pb.QueryTicketIdsRequest{
			Pool: &pb.Pool{
				Name: "default",
				TagPresentFilters: []*pb.TagPresentFilter{
					{
						Tag: userTag,
					},
				},
			},
		})
		if err != nil {
			return err
		}
		_, err = ticketClient.Recv()
		if err != io.EOF {
			c.Out = &api.CreateMatchmakingTicketResponse{
				Success: false,
				Error:   api.CreateMatchmakingTicketResponse_USER_ALREADY_HAS_ACTIVE_TICKET,
			}
			return nil
		}
		floatTicketMap[userTag] = float64(user.Elo)
		userTags = append(userTags, userTag)
	}
	// Create the ticket
	floatTicketMap["NumberOfUsers"] = float64(numberOfUsers)
	ticket, err := c.service.frontEndClient.CreateTicket(ctx, &pb.CreateTicketRequest{
		Ticket: &pb.Ticket{
			SearchFields: &pb.SearchFields{
				Tags:       append(arenaTags, userTags...),
				DoubleArgs: floatTicketMap,
			},
			Extensions: map[string]*anypb.Any{
				"Data": data,
			},
		},
	})
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	ticketId, err := strconv.Atoi(ticket.Id)
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
