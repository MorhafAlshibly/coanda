package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/matchmaking/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetMatchesCommand struct {
	service *Service
	In      *api.GetMatchesRequest
	Out     *api.GetMatchesResponse
}

func NewGetMatchesCommand(service *Service, in *api.GetMatchesRequest) *GetMatchesCommand {
	return &GetMatchesCommand{
		service: service,
		In:      in,
	}
}

func (c *GetMatchesCommand) Execute(ctx context.Context) error {
	matchLimit, matchOffset := conversion.PaginationToLimitOffset(c.In.Pagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	ticketLimit, ticketOffset := conversion.PaginationToLimitOffset(c.In.TicketPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	userLimit, userOffset := conversion.PaginationToLimitOffset(c.In.UserPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	arenaLimit, arenaOffset := conversion.PaginationToLimitOffset(c.In.ArenaPagination, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	statuses := []string{}
	for _, status := range c.In.Statuses {
		statuses = append(statuses, status.String())
	}
	matches, err := c.service.database.GetMatches(ctx, model.GetMatchesParams{
		Arena:           arenaRequestToArenaParams(c.In.Arena),
		MatchmakingUser: matchmakingUserRequestToMatchmakingUserParams(c.In.MatchmakingUser),
		Statuses:        statuses,
		Limit:           matchLimit,
		Offset:          matchOffset,
		TicketLimit:     ticketLimit,
		TicketOffset:    ticketOffset,
		UserLimit:       userLimit,
		UserOffset:      userOffset,
		ArenaLimit:      arenaLimit,
		ArenaOffset:     arenaOffset,
	})
	if err != nil {
		return err
	}
	apiMatches, err := unmarshalMatches(matches)
	if err != nil {
		return err
	}
	c.Out = &api.GetMatchesResponse{
		Success: true,
		Matches: apiMatches,
	}
	return nil
}
