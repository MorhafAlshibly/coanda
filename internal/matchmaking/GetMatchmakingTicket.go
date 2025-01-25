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
	"google.golang.org/protobuf/types/known/structpb"
	"open-match.dev/open-match/pkg/pb"
)

type GetMatchmakingTicketCommand struct {
	service *Service
	In      *api.GetMatchmakingTicketRequest
	Out     *api.GetMatchmakingTicketResponse
}

func NewGetMatchmakingTicketCommand(service *Service, in *api.GetMatchmakingTicketRequest) *GetMatchmakingTicketCommand {
	return &GetMatchmakingTicketCommand{
		service: service,
		In:      in,
	}
}

func (c *GetMatchmakingTicketCommand) Execute(ctx context.Context) error {
	mtErr := c.service.checkForMatchmakingTicketRequestError(c.In.MatchmakingTicket)
	if mtErr != nil {
		c.Out = &api.GetMatchmakingTicketResponse{
			Success: false,
			Error:   conversion.Enum(*mtErr, api.GetMatchmakingTicketResponse_Error_value, api.GetMatchmakingTicketResponse_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED),
		}
		return nil
	}
	// Make sure matchmaking user isnt nil
	if c.In.MatchmakingTicket.MatchmakingUser == nil {
		c.In.MatchmakingTicket.MatchmakingUser = &api.MatchmakingUserRequest{}
	}
	userLimit, userOffset := conversion.PaginationToLimitOffset(c.In.UserPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	arenaLimit, arenaOffset := conversion.PaginationToLimitOffset(c.In.ArenaPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	// Get matchmaking ticket
	var matchmakingTicket *pb.Ticket
	var err error
	if c.In.MatchmakingTicket.Id != nil {
		matchmakingTicket, err = c.service.frontEndClient.GetTicket(ctx, &pb.GetTicketRequest{
			TicketId: strconv.Itoa(int(*c.In.MatchmakingTicket.Id)),
		})
		if err != nil {
			return err
		}
	} else {
		if c.In.MatchmakingTicket.MatchmakingUser.ClientUserId == nil {
			user, err := c.service.database.GetMatchmakingUser(ctx, model.GetMatchmakingUserParams{
				Id: conversion.Uint64ToSqlNullInt64(c.In.MatchmakingTicket.MatchmakingUser.Id),
			})
			if err != nil {
				if err == sql.ErrNoRows {
					c.Out = &api.GetMatchmakingTicketResponse{
						Success: false,
						Error:   api.GetMatchmakingTicketResponse_NOT_FOUND,
					}
					return nil
				}
				return err
			}
			c.In.MatchmakingTicket.MatchmakingUser.ClientUserId = &user.ClientUserID
		}
		ticketClient, err := c.service.queryServiceClient.QueryTickets(ctx, &pb.QueryTicketsRequest{
			Pool: &pb.Pool{
				Name: "default",
				TagPresentFilters: []*pb.TagPresentFilter{
					{
						Tag: fmt.Sprintf("User_%d", c.In.MatchmakingTicket.MatchmakingUser.ClientUserId),
					},
				},
			},
		})
		if err != nil {
			return err
		}
		resp, err := ticketClient.Recv()
		if err != nil {
			if err == io.EOF {
				c.Out = &api.GetMatchmakingTicketResponse{
					Success: false,
					Error:   api.GetMatchmakingTicketResponse_NOT_FOUND,
				}
				return nil
			}
			return err
		}
		matchmakingTicket = resp.Tickets[0]
	}
	ticketUsers := make([]*api.MatchmakingUser, 0, userLimit)
	startIndex := userOffset
	endIndex := userOffset + userLimit
	for i := startIndex; i < endIndex; i++ {
		if int(i) >= len(matchmakingTicket.SearchFields.Tags) {
			break
		}
		tag := matchmakingTicket.SearchFields.Tags[i]
		// Check if tag is a user tag (starts with User_)
		if len(tag) < 5 || tag[:5] != "User_" {
			continue
		}
		clientUserID, err := strconv.Atoi(tag[5:])
		if err != nil {
			return err
		}
		user, err := c.service.database.GetMatchmakingUser(ctx, model.MatchmakingUserParams{
			ClientUserID: conversion.Int64ToSqlNullInt64(conversion.ValueToPointer(int64(clientUserID))),
		})
		if err != nil {
			return err
		}
		apiUser, err := unmarshalMatchmakingUser(user)
		if err != nil {
			return err
		}
		ticketUsers = append(ticketUsers, apiUser)
	}
	ticketArenas := make([]*api.Arena, 0, arenaLimit)
	startIndex = arenaOffset
	endIndex = arenaOffset + arenaLimit
	for i := startIndex; i < endIndex; i++ {
		if int(i) >= len(matchmakingTicket.SearchFields.Tags) {
			break
		}
		tag := matchmakingTicket.SearchFields.Tags[i]
		// Check if tag is an arena tag (starts with Arena_)
		if len(tag) < 6 || tag[:6] != "Arena_" {
			continue
		}
		arenaID, err := strconv.Atoi(tag[6:])
		if err != nil {
			return err
		}
		arena, err := c.service.database.GetArena(ctx, model.ArenaParams{
			ID: conversion.Int64ToSqlNullInt64(conversion.ValueToPointer(int64(arenaID))),
		})
		if err != nil {
			return err
		}
		apiArena, err := unmarshalArena(arena)
		if err != nil {
			return err
		}
		ticketArenas = append(ticketArenas, apiArena)
	}
	status := api.MatchmakingTicket_PENDING
	if matchmakingTicket.Assignment != nil {
		status = api.MatchmakingTicket_MATCHED
	}
	data := new(structpb.Struct)
	if matchmakingTicket.Extensions["Data"] != nil {
		err = matchmakingTicket.Extensions["Data"].UnmarshalTo(data)
		if err != nil {
			return err
		}
	}
	out := &api.MatchmakingTicket{
		Id:               conversion.ValueToPointer(uint64(matchmakingTicket.TicketId)),
		MatchmakingUsers: ticketUsers,
		Arenas:           ticketArenas,
		MatchId:          conversion.ValueToPointer(matchmakingTicket.Assignment.Connection),
		Status:           status,
		Data:             data,
		CreatedAt:        matchmakingTicket.CreateTime,
		UpdatedAt:        matchmakingTicket.CreateTime,
	}
	c.Out = &api.GetMatchmakingTicketResponse{
		Success:           true,
		MatchmakingTicket: out,
		Error:             api.GetMatchmakingTicketResponse_NONE,
	}
	return nil
}
