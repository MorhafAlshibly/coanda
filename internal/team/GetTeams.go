package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
)

type GetTeamsCommand struct {
	service *Service
	In      *api.Pagination
	Out     *api.GetTeamsResponse
}

func NewGetTeamsCommand(service *Service, in *api.Pagination) *GetTeamsCommand {
	return &GetTeamsCommand{
		service: service,
		In:      in,
	}
}

func (c *GetTeamsCommand) Execute(ctx context.Context) error {
	limit, offset := conversion.PaginationToLimitOffset(c.In, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	teams, err := c.service.database.GetTeams(ctx, model.GetTeamsParams{
		Limit:  limit,
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
	c.Out = &api.GetTeamsResponse{
		Success: true,
		Teams:   outs,
	}
	return nil
}
