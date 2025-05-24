package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetMatchCommand struct {
	service *Service
	In      *api.GetMatchRequest
	Out     *api.GetMatchResponse
}

func NewGetMatchCommand(service *Service, in *api.GetMatchRequest) *GetMatchCommand {
	return &GetMatchCommand{
		service: service,
		In:      in,
	}
}

func (c *GetMatchCommand) Execute(ctx context.Context) error {
	mmErr := c.service.checkForMatchRequestError(c.In.Match)
	if mmErr != nil {
		c.Out = &api.GetMatchResponse{
			Success: false,
			Error:   conversion.Enum(*mmErr, api.GetMatchResponse_Error_value, api.GetMatchResponse_MATCH_ID_OR_MATCHMAKING_TICKET_REQUIRED),
		}
		return nil
	}
	ticketLimit, ticketOffset := conversion.PaginationToLimitOffset(c.In.TicketPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	userLimit, userOffset := conversion.PaginationToLimitOffset(c.In.UserPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	arenaLimit, arenaOffset := conversion.PaginationToLimitOffset(c.In.ArenaPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	match, err := c.service.database.GetMatch(ctx, model.GetMatchParams{
		Match:        matchRequestToMatchParams(c.In.Match),
		TicketLimit:  ticketLimit,
		TicketOffset: ticketOffset,
		UserLimit:    userLimit,
		UserOffset:   userOffset,
		ArenaLimit:   arenaLimit,
		ArenaOffset:  arenaOffset,
	})
	if err != nil {
		return err
	}
	if len(match) == 0 {
		c.Out = &api.GetMatchResponse{
			Success: false,
			Error:   api.GetMatchResponse_NOT_FOUND,
		}
		return nil
	}
	apiMatch, err := unmarshalMatch(match)
	if err != nil {
		return err
	}
	c.Out = &api.GetMatchResponse{
		Success: true,
		Match:   apiMatch,
		Error:   api.GetMatchResponse_NONE,
	}
	return nil
}
