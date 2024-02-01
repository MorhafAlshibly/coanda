package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/validation"
)

type SearchTeamsCommand struct {
	service *Service
	In      *api.SearchTeamsRequest
	Out     *api.SearchTeamsResponse
}

func NewSearchTeamsCommand(service *Service, in *api.SearchTeamsRequest) *SearchTeamsCommand {
	return &SearchTeamsCommand{
		service: service,
		In:      in,
	}
}

func (c *SearchTeamsCommand) Execute(ctx context.Context) error {
	if len(c.In.Query) < int(c.service.minTeamNameLength) {
		c.Out = &api.SearchTeamsResponse{
			Success: false,
			Error:   api.SearchTeamsResponse_QUERY_TOO_SHORT,
		}
		return nil
	}
	if len(c.In.Query) > int(c.service.maxTeamNameLength) {
		c.Out = &api.SearchTeamsResponse{
			Success: false,
			Error:   api.SearchTeamsResponse_QUERY_TOO_LONG,
		}
		return nil
	}
	if c.In.Pagination == nil {
		c.In.Pagination = &api.GetTeamsRequest{}
	}
	max := validation.ValidateMaxPageLength(c.In.Pagination.Max, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	offset := conversion.PageToOffset(c.In.Pagination.Page, max)
	teams, err := c.service.database.SearchTeams(ctx, model.SearchTeamsParams{
		Query:  c.In.Query,
		Limit:  int32(max),
		Offset: offset,
	})
	if err != nil {
		return err
	}
	outs := make([]*api.Team, len(teams))
	for i, team := range teams {
		outs[i], err = UnmarshalTeam(team)
		if err != nil {
			return err
		}
	}
	c.Out = &api.SearchTeamsResponse{
		Success: true,
		Teams:   outs,
		Error:   api.SearchTeamsResponse_NONE,
	}
	return nil

}
