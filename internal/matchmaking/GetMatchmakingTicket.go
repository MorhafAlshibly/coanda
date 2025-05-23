package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
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
			Error:   conversion.Enum(*mtErr, api.GetMatchmakingTicketResponse_Error_value, api.GetMatchmakingTicketResponse_MATCHMAKING_TICKET_ID_OR_MATCHMAKING_USER_REQUIRED),
		}
		return nil
	}
	userLimit, userOffset := conversion.PaginationToLimitOffset(c.In.UserPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	arenaLimit, arenaOffset := conversion.PaginationToLimitOffset(c.In.ArenaPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	matchmakingTicket, err := c.service.database.GetMatchmakingTicket(ctx, model.GetMatchmakingTicketParams{
		MatchmakingTicket: matchmakingTicketRequestToMatchmakingTicketParams(c.In.MatchmakingTicket),
		UserLimit:         userLimit,
		UserOffset:        userOffset,
		ArenaLimit:        arenaLimit,
		ArenaOffset:       arenaOffset,
	})
	if err != nil {
		return err
	}
	if len(matchmakingTicket) == 0 {
		c.Out = &api.GetMatchmakingTicketResponse{
			Success: false,
			Error:   api.GetMatchmakingTicketResponse_NOT_FOUND,
		}
		return nil
	}
	apiMatchmakingTicket, err := unmarshalMatchmakingTicket(matchmakingTicket)
	if err != nil {
		return err
	}
	c.Out = &api.GetMatchmakingTicketResponse{
		Success:           true,
		MatchmakingTicket: apiMatchmakingTicket,
		Error:             api.GetMatchmakingTicketResponse_NONE,
	}
	return nil
}
