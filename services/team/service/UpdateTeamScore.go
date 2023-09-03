package service

import (
	"context"

	"github.com/MorhafAlshibly/coanda/services/team"
	"github.com/MorhafAlshibly/coanda/services/team/schema"
	"go.mongodb.org/mongo-driver/bson"
)

type UpdateTeamScoreCommand struct {
	Service *team.TeamService
	In      *schema.UpdateTeamScoreRequest
	Out     *schema.Team
}

func (c *UpdateTeamScoreCommand) Execute(ctx context.Context) error {
	filter, err := team.GetFilter(c.In.Team)
	if err != nil {
		return err
	}
	_, err = c.Service.Db.UpdateOne(ctx, filter, bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "score", Value: c.In.ScoreOffset},
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
