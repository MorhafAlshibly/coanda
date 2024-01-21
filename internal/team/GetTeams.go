package team

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/MorhafAlshibly/coanda/internal/team/model"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"github.com/MorhafAlshibly/coanda/pkg/validation"
)

type GetTeamsCommand struct {
	service *Service
	In      *api.GetTeamsRequest
	Out     *api.GetTeamsResponse
}

func NewGetTeamsCommand(service *Service, in *api.GetTeamsRequest) *GetTeamsCommand {
	return &GetTeamsCommand{
		service: service,
		In:      in,
	}
}

func (c *GetTeamsCommand) Execute(ctx context.Context) error {
	max := validation.ValidateMaxPageLength(c.In.Max, c.service.defaultMaxPageLength, c.service.maxMaxPageLength)
	offset := conversion.PageToOffset(c.In.Page, max)
	teams, err := c.service.database.GetTeams(ctx, model.GetTeamsParams{
		Limit:  int32(max),
		Offset: offset,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Out = &api.GetTeamsResponse{
				Success: false,
				Teams:   []*api.Team{},
			}
			return nil
		}
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
