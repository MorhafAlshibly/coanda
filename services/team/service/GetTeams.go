package service

import (
	"context"

	"github.com/MorhafAlshibly/coanda/services/team/schema"
)

type GetTeamsCommand struct {
	Service *TeamService
	In      *schema.GetTeamsRequest
	Out     *schema.Teams
}

func (c *GetTeamsCommand) Execute(ctx context.Context) error {
	cursor, err := c.Service.Db.Aggregate(ctx, c.Service.Pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	c.Out, err = ToTeams(ctx, cursor, c.In.Page, c.In.Max)
	if err != nil {
		return err
	}
	return nil
}
