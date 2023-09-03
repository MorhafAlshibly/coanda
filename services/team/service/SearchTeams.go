package service

import (
	"context"
	"errors"

	"github.com/MorhafAlshibly/coanda/services/team/schema"
	"go.mongodb.org/mongo-driver/bson"
)

type SearchTeamsCommand struct {
	Service *TeamService
	In      *schema.SearchTeamsRequest
	Out     *schema.Teams
}

func (c *SearchTeamsCommand) Execute(ctx context.Context) error {
	searchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "name", Value: bson.D{
				{Key: "$regex", Value: c.In.Query},
				{Key: "$options", Value: "i"},
			}},
		}},
	}
	if len(c.In.Query) < c.Service.MinTeamNameLength {
		return errors.New("query too short")
	}
	cursor, err := c.Service.Db.Aggregate(ctx, append(c.Service.Pipeline, searchStage))
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
