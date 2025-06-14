package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type PollMatchmakingTicketCommand struct {
	service *Service
	In      *api.GetMatchmakingTicketRequest
	Out     *api.PollMatchmakingTicketResponse
}

func NewPollMatchmakingTicketCommand(service *Service, in *api.GetMatchmakingTicketRequest) *PollMatchmakingTicketCommand {
	return &PollMatchmakingTicketCommand{
		service: service,
		In:      in,
	}
}

func (c *PollMatchmakingTicketCommand) Execute(ctx context.Context) error {
	mtErr := c.service.checkForMatchmakingTicketRequestError(c.In.MatchmakingTicket)
	if mtErr != nil {
		c.Out = &api.PollMatchmakingTicketResponse{
			Success: false,
			Error:   conversion.Enum(*mtErr, api.PollMatchmakingTicketResponse_Error_value, api.PollMatchmakingTicketResponse_MATCHMAKING_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED),
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
	params := matchmakingTicketRequestToMatchmakingTicketParams(c.In.MatchmakingTicket)
	_, err = qtx.PollMatchmakingTicket(ctx, model.PollMatchmakingTicketParams{
		MatchmakingTicket: params,
		ExpiryTimeWindow:  c.service.expiryTimeWindow,
	})
	if err != nil {
		return err
	}
	userLimit, userOffset := conversion.PaginationToLimitOffset(c.In.UserPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	arenaLimit, arenaOffset := conversion.PaginationToLimitOffset(c.In.ArenaPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	matchmakingTicket, err := qtx.GetMatchmakingTicket(ctx, model.GetMatchmakingTicketParams{
		MatchmakingTicket: params,
		UserLimit:         userLimit,
		UserOffset:        userOffset,
		ArenaLimit:        arenaLimit,
		ArenaOffset:       arenaOffset,
	})
	if err != nil {
		return err
	}
	if len(matchmakingTicket) == 0 {
		c.Out = &api.PollMatchmakingTicketResponse{
			Success: false,
			Error:   api.PollMatchmakingTicketResponse_NOT_FOUND,
		}
		return nil
	}
	if matchmakingTicket[0].Status == "EXPIRED" {
		c.Out = &api.PollMatchmakingTicketResponse{
			Success: false,
			Error:   api.PollMatchmakingTicketResponse_ALREADY_EXPIRED,
		}
		return nil
	}
	if matchmakingTicket[0].Status == "MATCHED" {
		c.Out = &api.PollMatchmakingTicketResponse{
			Success: false,
			Error:   api.PollMatchmakingTicketResponse_ALREADY_MATCHED,
		}
		return nil
	}
	if matchmakingTicket[0].Status == "ENDED" {
		c.Out = &api.PollMatchmakingTicketResponse{
			Success: false,
			Error:   api.PollMatchmakingTicketResponse_ALREADY_ENDED,
		}
		return nil
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	apiMatchmakingTicket, err := unmarshalMatchmakingTicket(matchmakingTicket)
	if err != nil {
		return err
	}
	c.Out = &api.PollMatchmakingTicketResponse{
		Success:           true,
		MatchmakingTicket: apiMatchmakingTicket,
		Error:             api.PollMatchmakingTicketResponse_NONE,
	}
	return nil
}
