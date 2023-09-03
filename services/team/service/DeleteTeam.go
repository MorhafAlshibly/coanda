package service

import (
	"context"

	"github.com/MorhafAlshibly/coanda/services/team/schema"
)

type DeleteTeamCommand struct {
	Service *TeamService
	In      *schema.DeleteTeamRequest
	Out     *schema.Team
}

func (c *DeleteTeamCommand) Execute(ctx context.Context) error {
	filter, err := GetFilter(c.In.Team)
	if err != nil {
		return err
	}
	c.Out, err = c.Service.GetTeam(ctx, c.In.Team)
	if err != nil {
		return err
	}
	_, err = c.Service.Db.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
