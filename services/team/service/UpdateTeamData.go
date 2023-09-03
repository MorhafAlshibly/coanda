package service

import (
	"context"

	"github.com/MorhafAlshibly/coanda/services/team"
	"github.com/MorhafAlshibly/coanda/services/team/schema"
	"go.mongodb.org/mongo-driver/bson"
)

type UpdateTeamDataCommand struct {
	Service *team.TeamService
	In      *schema.UpdateTeamDataRequest
	Out     *schema.Team
}

func (c *UpdateTeamDataCommand) Execute(ctx context.Context) error {
	filter, err := team.GetFilter(c.In.Team)
	if err != nil {
		return err
	}
	_, err = c.Service.Db.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "data", Value: c.In.Data},
		}},
	})
	if err != nil {
		return err
	}
	c.Out, err = c.Service.GetTeam(ctx, c.In.Team)
	if err != nil {
		return err
	}
	return nil
}
