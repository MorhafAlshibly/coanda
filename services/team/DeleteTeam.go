package team

import (
	"context"

	"github.com/MorhafAlshibly/coanda/services/team/schema"
)

type DeleteTeamCommand struct {
	service *TeamService
	In      *schema.DeleteTeamRequest
	Out     *schema.Team
}

func NewDeleteTeamCommand(service *TeamService, in *schema.DeleteTeamRequest) *DeleteTeamCommand {
	return &DeleteTeamCommand{
		service: service,
		In:      in,
	}
}

func (c *DeleteTeamCommand) Execute(ctx context.Context) error {
	filter, err := getFilter(c.In.Team)
	if err != nil {
		return err
	}
	c.Out, err = c.service.GetTeam(ctx, c.In.Team)
	if err != nil {
		return err
	}
	_, err = c.service.Db.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
