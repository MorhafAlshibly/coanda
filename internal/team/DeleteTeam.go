package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type DeleteTeamCommand struct {
	service *Service
	In      *api.DeleteTeamRequest
	Out     *api.GetTeamResponse
}

func NewDeleteTeamCommand(service *Service, in *api.DeleteTeamRequest) *DeleteTeamCommand {
	return &DeleteTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteTeamCommand) Execute(ctx context.Context) error {
	filter, err := getFilter(c.In.Team)
	if err != nil {
		c.Out = &api.GetTeamResponse{
			Success: false,
			Team:    nil,
			Error:   api.GetTeamResponse_INVALID,
		}
		return nil
	}
	c.Out, err = c.service.GetTeam(ctx, c.In.Team)
	if err != nil {
		return err
	}
	_, err = c.service.db.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
